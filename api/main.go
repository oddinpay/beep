package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"go.jetify.com/sse"
	"go.jetify.com/typeid/v2"
)

const (
	Host = "0.0.0.0"
	Port = "8976"

	// HTTP methods
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodPatch   = "PATCH"
	MethodOptions = "OPTIONS"

	// Content types
	ContentTypeJSON        = "application/json"
	ContentTypeEventStream = "text/event-stream"

	// HTTP State codes
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusMethodNotAllowed    = 405
	StatusMultipleChoices     = 300

	// Common headers
	HeaderContentType  = "Content-Type"
	HeaderCacheControl = "Cache-Control"
	HeaderConnection   = "Connection"
	HeaderAllowOrigin  = "Access-Control-Allow-Origin"
	HeaderAllowMethods = "Access-Control-Allow-Methods"
	HeaderAllowHeaders = "Access-Control-Allow-Headers"

	defaultTimeout = 20 * time.Second
	minutes90d     = 90
)

// ----------- DB / CACHE CONNECTIONS -----------

var (
	jwt              = os.Getenv("NATS_JWT")
	seed             = os.Getenv("NATS_SEED")
	serverURL        = os.Getenv("NATS_URL")
	userAgent        = os.Getenv("USER_AGENT")
	probeManagerOnce sync.Once
	monitorStartTime = time.Now().UTC().Truncate(24 * time.Hour)
	hr               = HealthResponse{Down: "down", Up: "up", Warn: "warn"}
	nc               *nats.Conn
	err              error
	wg               sync.WaitGroup
	js               jetstream.JetStream
	kv               jetstream.KeyValue
)

// -------------------- GLOBAL SLA MAP --------------------

var slaTrackers = struct {
	sync.Mutex
	m map[string]*SlidingSLA
}{m: make(map[string]*SlidingSLA)}

var defaultReqs = func() []HttpRequest {
	raw := []HttpRequest{
		{Name: "www.oddinpay.com", Protocol: "https", Host: "www.oddinpay.com", Interval: 10 * time.Second},
	}

	// for i := 1; i <= 2; i++ {
	// 	raw = append(raw, HttpRequest{Name: fmt.Sprintf("DNS %d", i), Protocol: "dns", Host: "www.oddinpay.com", Interval: 10 * time.Second})
	// 	raw = append(raw, HttpRequest{Name: fmt.Sprintf("HTTPS %d", i), Protocol: "https", Host: "www.oddinpay.com", Interval: 10 * time.Second})
	// }

	out := make([]HttpRequest, 0, len(raw))
	counts := make(map[string]int)

	for _, r := range raw {
		name := r.Name
		counts[name]++
		if counts[name] > 1 {
			r.Name = fmt.Sprintf("%s-%d", name, counts[name])

		}
		out = append(out, r)
	}
	return out
}()

// -------------------- MODELS --------------------

type HttpRequest struct {
	Host     string        `json:"host,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
	Name     string        `json:"name,omitempty"`
	Username string        `json:"username,omitempty"`
	Password string        `json:"password,omitempty"`
}

type HealthResponse struct {
	Down string `json:"down"`
	Up   string `json:"up"`
	Warn string `json:"warn"`
}

type ProbeResult struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Protocol    string   `json:"protocol,omitempty"`
	State       []string `json:"state,omitempty"`
	Description string   `json:"description,omitempty"`
	Date        []string `json:"date,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
}

type ProbeResponse struct {
	Index   int           `json:"index"`
	Payload StatusPayload `json:"payload"`
}

type StatusPayload struct {
	Probe ProbeResult    `json:"probe"`
	SLA   map[string]any `json:"sla"`
}

type ErrorResponse struct {
	State   []string `json:"state"`
	Message string   `json:"message"`
}

type bucket struct{ totalSec, downSec int64 }

type SlidingSLA struct {
	Target        float64
	buckets       []bucket
	idx           int
	currentMinute time.Time
	lastUpdate    time.Time
	mu            sync.Mutex
}

// -------------------- RECOVERY --------------------

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Recovered from panic: %v\n", rec)
				if w.Header().Get(HeaderContentType) == "" {
					w.Header().Set(HeaderContentType, ContentTypeJSON)
				}
				w.WriteHeader(StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					State:   []string{"error"},
					Message: "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func parseDurationToSecs(s string) int64 {
	var total int64
	parts := strings.FieldsSeq(s)
	for part := range parts {
		var val int64
		if strings.HasSuffix(part, "d") {
			fmt.Sscanf(part, "%dd", &val)
			total += val * 86400
		} else if strings.HasSuffix(part, "h") {
			fmt.Sscanf(part, "%dh", &val)
			total += val * 3600
		} else if strings.HasSuffix(part, "m") {
			fmt.Sscanf(part, "%dm", &val)
			total += val * 60
		} else if strings.HasSuffix(part, "s") {
			fmt.Sscanf(part, "%ds", &val)
			total += val
		}
	}
	return total
}

func formatDurationFull(seconds int64) string {
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

func getRecentDates() []string {
	return []string{time.Now().UTC().Format("02/01/2006")}
}

// -------------------- BROADCAST HUB --------------------

type Hub struct {
	sync.RWMutex
	clients map[chan map[string]StatusPayload]struct{}
	cache   map[string]StatusPayload
}

var globalHub = &Hub{
	clients: make(map[chan map[string]StatusPayload]struct{}),
	cache:   make(map[string]StatusPayload),
}

func (h *Hub) Broadcast(update map[string]StatusPayload) {
	h.Lock()
	defer h.Unlock()

	maps.Copy(h.cache, update)

	for clientChan := range h.clients {
		select {
		case clientChan <- update:
		default:
		}
	}
}

func monitorId() string {
	monitorId := typeid.MustGenerate("monitor")
	return monitorId.String()
}

func slaId() string {
	slaId := typeid.MustGenerate("sla")
	return slaId.String()
}

// -------------------- PROBES --------------------

func probeHTTP(re HttpRequest) ProbeResult {

	url := fmt.Sprintf("%s://%s", re.Protocol, re.Host)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Failed to create HTTP request", "error", err)
	}

	if userAgent == "" {
		userAgent = "BeepMonitor/1.0"
	}

	r.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return ProbeResult{
			Id:          "",
			Name:        re.Name,
			Protocol:    strings.ToUpper(re.Protocol),
			Description: fmt.Sprintf("%s - %s", re.Host, err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Down},
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < StatusOK || resp.StatusCode >= StatusBadRequest {
		return ProbeResult{
			Id:          "",
			Name:        re.Name,
			Protocol:    strings.ToUpper(re.Protocol),
			Description: fmt.Sprintf("%s - %d", re.Host, resp.StatusCode),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Down},
		}
	}
	return ProbeResult{
		Id:          "",
		Name:        re.Name,
		Protocol:    strings.ToUpper(re.Protocol),
		Description: fmt.Sprintf("%s - %d", re.Host, resp.StatusCode),
		Timestamp:   time.Now().Format("15:04:05.000"),
		Date:        getRecentDates(),
		State:       []string{hr.Up},
	}
}

func probeTCP(req HttpRequest) ProbeResult {
	conn, err := net.DialTimeout("tcp", req.Host, defaultTimeout)

	if err != nil {
		return ProbeResult{
			Id:          "",
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: err.Error(),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Down},
		}
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping\n"))
	if err != nil {
		return ProbeResult{
			Id:          "",
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: "write failed: " + err.Error(),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Down},
		}
	}

	buf := make([]byte, 64)
	_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ProbeResult{
			Id:          "",
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: "no response after connect",
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Up},
		}
	}

	return ProbeResult{
		Id:          "",
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		Description: fmt.Sprintf("response received %s", strings.TrimSpace(string(buf[:n]))),
		Timestamp:   time.Now().Format("15:04:05.000"),
		Date:        getRecentDates(),
		State:       []string{hr.Up},
	}
}

func probeDNS(req HttpRequest) ProbeResult {

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if net.ParseIP(req.Host) != nil {
		return ProbeResult{
			Id:          "",
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: "Input is already an IP, DNS lookup skipped",
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Warn},
		}
	}

	addrs, err := net.DefaultResolver.LookupHost(ctx, req.Host)
	if err != nil {
		return ProbeResult{
			Id:          "",
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: fmt.Sprintf("DNS error: %s", err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Down},
		}
	}

	return ProbeResult{
		Id:          "",
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		Description: fmt.Sprintf("resolved %v", addrs),
		Timestamp:   time.Now().Format("15:04:05.000"),
		Date:        getRecentDates(),
		State:       []string{hr.Up, "up"},
	}
}

// -------------------- 90-DAY SLA --------------------

func NewSlidingSLA(target float64) *SlidingSLA {
	now := time.Now()
	return &SlidingSLA{
		Target:        target,
		buckets:       make([]bucket, minutes90d),
		currentMinute: now.Truncate(24 * time.Hour),
		lastUpdate:    now,
	}
}

func (s *SlidingSLA) rotateTo(now time.Time) {
	minNow := now.Truncate(24 * time.Hour)
	if !minNow.After(s.currentMinute) {
		return
	}
	steps := int(minNow.Sub(s.currentMinute) / 24 * time.Hour)
	if steps > minutes90d {
		for i := range s.buckets {
			s.buckets[i] = bucket{}
		}
		s.idx = 0
		s.currentMinute = minNow
		return
	}
	for range steps {
		s.idx++
		if s.idx >= len(s.buckets) {
			s.idx = 0
		}
		s.buckets[s.idx] = bucket{}
	}
	s.currentMinute = minNow
}

func (s *SlidingSLA) Tick(isDown bool, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.rotateTo(now)

	inc := int64(interval.Round(time.Second).Seconds())

	s.buckets[s.idx].totalSec += inc
	if isDown {
		s.buckets[s.idx].downSec += inc
	}
	s.lastUpdate = now
}

func (s *SlidingSLA) Snapshot() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	var total, down int64
	for _, b := range s.buckets {
		total += b.totalSec
		down += b.downSec
	}

	if total <= 0 {
		return map[string]any{
			"id":                 "",
			"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
			"uptime90":           "100.000%",
			"up_time_seconds":    formatDurationFull(0),
			"down_time_seconds":  formatDurationFull(0),
			"total_time_seconds": formatDurationFull(0),
			"sla_breached":       false,
		}
	}

	availability := 1.0 - (float64(down) / float64(total))
	percent := availability * 100

	uptimeStr := fmt.Sprintf("%.3f%%", percent)
	if down > 0 && uptimeStr == "100.000%" {
		uptimeStr = "99.999%"
	}

	breached := (s.Target >= 1.0 && down > 0) || (availability < s.Target)
	up := total - down

	return map[string]any{
		"id":                 "",
		"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
		"uptime90":           uptimeStr,
		"up_time_seconds":    formatDurationFull(up),
		"down_time_seconds":  formatDurationFull(down),
		"total_time_seconds": formatDurationFull(total),
		"sla_breached":       breached,
	}
}

func (s *SlidingSLA) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.buckets {
		s.buckets[i] = bucket{}
	}
	s.idx = 0
	s.currentMinute = time.Now().Truncate(time.Minute)
	s.lastUpdate = time.Now()
}

func startProbeManager(ctx context.Context, wg *sync.WaitGroup) {
	probeManagerOnce.Do(func() {
		slog.Info("Starting probe manager...")

		for _, target := range defaultReqs {

			t := target

			var probeFn func(HttpRequest) ProbeResult
			switch strings.ToLower(strings.TrimSpace(t.Protocol)) {
			case "tcp":
				probeFn = probeTCP
			case "http", "https":
				probeFn = probeHTTP
			case "dns":
				probeFn = probeDNS
			}

			if probeFn == nil {
				slog.Warn("Unsupported protocol", "protocol", t.Protocol)
				continue
			}

			wg.Add(1)

			interval := t.Interval
			if interval <= 0 {
				interval = 1 * time.Second
			}

			go func(req HttpRequest, fn func(HttpRequest) ProbeResult, iv time.Duration) {
				defer wg.Done()

				ticker := time.NewTicker(iv)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						slog.Info("Stopping probe worker", "name", req.Name)
						return

					case <-ticker.C:
						ctx, cancel := context.WithTimeout(ctx, defaultTimeout)

						res := fn(req)

						slaTrackers.Lock()
						tracker := slaTrackers.m[req.Name]
						if tracker == nil {
							tracker = NewSlidingSLA(1.0)
							slaTrackers.m[req.Name] = tracker
						}
						slaTrackers.Unlock()

						isDown := len(res.State) > 0 && strings.ToLower(res.State[0]) == hr.Down
						tracker.Tick(isDown, interval)

						payload := StatusPayload{
							Probe: res,
							SLA:   tracker.Snapshot(),
						}

						publishToNATS(ctx, req.Name, &payload, tracker)

						// Broadcast update
						globalHub.Broadcast(map[string]StatusPayload{req.Name: payload})

						cancel()

					}
				}

			}(t, probeFn, interval)
		}
	})

}

// -------------------- SSE HANDLER --------------------

func Sse(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// SSE headers
	w.Header().Set(HeaderAllowOrigin, "*")
	w.Header().Set(HeaderCacheControl, "no-cache")
	w.Header().Set(HeaderConnection, "keep-alive")
	w.Header().Set(HeaderContentType, ContentTypeEventStream)

	conn, err := sse.Upgrade(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), StatusInternalServerError)
		return
	}
	defer conn.Close()

	clientChan := make(chan map[string]StatusPayload, 100)

	globalHub.Lock()
	globalHub.clients[clientChan] = struct{}{}

	initialSnap := make(map[string]StatusPayload)
	maps.Copy(initialSnap, globalHub.cache)

	globalHub.Unlock()

	defer func() {
		globalHub.Lock()
		delete(globalHub.clients, clientChan)
		globalHub.Unlock()
	}()

	if len(initialSnap) > 0 {
		if err := sendUpdateToConn(ctx, conn, initialSnap); err != nil {
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-clientChan:

			if err := sendUpdateToConn(ctx, conn, update); err != nil {
				return
			}

		}
	}
}

func sendUpdateToConn(ctx context.Context, conn *sse.Conn, update map[string]StatusPayload) error {
	for name, payload := range update {
		idx := -1
		for i, r := range defaultReqs {
			if r.Name == name {
				idx = i
				break
			}
		}

		out := map[string]any{
			"index": idx,
			"payload": map[string]any{
				"probe": payload.Probe,
				"sla":   payload.SLA,
			},
		}
		if err := conn.SendData(ctx, out); err != nil {
			return err
		}
	}
	return nil
}

// -------------------- STATE REQUEST HANDLER --------------------

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(HeaderContentType, ContentTypeJSON)

	var hasMonitors bool
	var miniMonitors bool
	if len(defaultReqs) == 0 {
		hasMonitors = false
	} else {
		hasMonitors = true
	}

	if len(defaultReqs) > 3 {
		miniMonitors = true
	} else {
		miniMonitors = false
	}

	response := map[string]bool{
		"monitors":     hasMonitors,
		"miniMonitors": miniMonitors,
	}

	respJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJSON)

}

// -------------------- SLA RESET HANDLER --------------------

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != MethodGet {
		http.Error(w, "Method not allowed", StatusMethodNotAllowed)
		return
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)

	name := r.URL.Query().Get("name")
	empty := r.URL.Query().Get("empty") == "true"

	slaTrackers.Lock()
	if name != "" {
		if tracker, ok := slaTrackers.m[name]; ok {
			tracker.Reset()
		}
	} else {
		for _, tracker := range slaTrackers.m {
			tracker.Reset()
		}
	}
	slaTrackers.Unlock()

	if empty {
		w.WriteHeader(StatusNoContent)
	}

	w.WriteHeader(StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"sla_reset": true,
		"probe":     name,
	})
}

func CreatePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != MethodPost {
		http.Error(w, "Method not allowed", StatusMethodNotAllowed)
		return
	}
}

func publishToNATS(ctx context.Context, name string, payload *StatusPayload, s *SlidingSLA) {
	if nc.Status() != nats.CONNECTED {
		slog.Error("NATS not connected")
		return
	}

	now := time.Now().UTC()

	// 2-minute block

	// intervalBlock := (now.Minute() / 2) * 2
	// todayUTC := fmt.Sprintf("%s %02d:%02d", now.Format("02/01/2006"), now.Hour(), intervalBlock)

	// Daily block
	todayUTC := now.Format("02/01/2006")

	currentStatus := hr.Warn
	if len(payload.Probe.State) > 0 {
		currentStatus = payload.Probe.State[0]
	}

	for range 3 {
		entry, getErr := kv.Get(ctx, name)
		var revision uint64 = 0
		var oldPayload StatusPayload

		if getErr == nil {
			revision = entry.Revision()
			gr, err := gzip.NewReader(bytes.NewReader(entry.Value()))
			if err == nil {
				decomp, _ := io.ReadAll(gr)
				var wrapped map[string]any
				json.Unmarshal(decomp, &wrapped)
				if p, ok := wrapped["payload"].(map[string]any); ok {

					existingProbeID, _ := p["probe"].(map[string]any)["id"].(string)
					existingSlaID, _ := p["sla"].(map[string]any)["id"].(string)

					if existingProbeID != "" {
						payload.Probe.Id = existingProbeID
					}
					if existingSlaID != "" {
						payload.SLA["id"] = existingSlaID
					}

					pBytes, _ := json.Marshal(p)
					json.Unmarshal(pBytes, &oldPayload)
				}
				gr.Close()
			}
		}

		if payload.Probe.Id == "" {
			payload.Probe.Id = monitorId()
		}
		if payload.SLA["id"] == nil || payload.SLA["id"] == "" {
			payload.SLA["id"] = slaId()
		}

		if getErr == nil && len(oldPayload.Probe.Date) > 0 {
			if oldPayload.Probe.Date[0] == todayUTC {
				payload.SLA["history"] = oldPayload.SLA["history"]
				payload.Probe.Date = oldPayload.Probe.Date
				payload.Probe.State = oldPayload.Probe.State

				if len(payload.Probe.State) > 0 {
					payload.Probe.State[0] = currentStatus
				} else {
					payload.Probe.State = []string{currentStatus}
				}

				if h, ok := payload.SLA["history"].([]any); ok && len(h) > 0 {
					h[0] = map[string]any{
						"sla_breached":       payload.SLA["sla_breached"],
						"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
						"total_time_seconds": payload.SLA["total_time_seconds"],
						"up_time_seconds":    payload.SLA["up_time_seconds"],
						"down_time_seconds":  payload.SLA["down_time_seconds"],
						"uptime90":           payload.SLA["uptime90"],
					}
				}
			} else {
				s.Reset()
				freshSLA := s.Snapshot()
				newSnapshot := map[string]any{
					"sla_breached":       freshSLA["sla_breached"],
					"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
					"total_time_seconds": freshSLA["total_time_seconds"],
					"up_time_seconds":    freshSLA["up_time_seconds"],
					"down_time_seconds":  freshSLA["down_time_seconds"],
					"uptime90":           freshSLA["uptime90"],
				}
				if oldHist, ok := oldPayload.SLA["history"].([]any); ok {
					payload.SLA["history"] = append([]any{newSnapshot}, oldHist...)
				}
				payload.Probe.Date = append([]string{todayUTC}, oldPayload.Probe.Date...)
				payload.Probe.State = append([]string{currentStatus}, oldPayload.Probe.State...)
			}
		} else {
			payload.SLA["history"] = []any{map[string]any{
				"sla_breached":       payload.SLA["sla_breached"],
				"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
				"total_time_seconds": payload.SLA["total_time_seconds"],
				"up_time_seconds":    payload.SLA["up_time_seconds"],
				"down_time_seconds":  payload.SLA["down_time_seconds"],
				"uptime90":           payload.SLA["uptime90"],
			}}
			payload.Probe.Date = []string{todayUTC}
		}

		payload.Probe.State = capSlice(payload.Probe.State, 90)
		payload.Probe.Date = capSlice(payload.Probe.Date, 90)
		if h, ok := payload.SLA["history"].([]any); ok {
			payload.SLA["history"] = capSlice(h, 90)
		}

		var rootTotal, rootDown int64
		if h, ok := payload.SLA["history"].([]any); ok {
			for _, hEntry := range h {
				if m, ok := hEntry.(map[string]any); ok {
					rootTotal += parseDurationToSecs(m["total_time_seconds"].(string))
					rootDown += parseDurationToSecs(m["down_time_seconds"].(string))
				}
			}
		}

		rootUp := rootTotal - rootDown
		rootAvail := 1.0
		if rootTotal > 0 {
			rootAvail = 1.0 - (float64(rootDown) / float64(rootTotal))
		}
		payload.SLA["total_time_seconds"] = formatDurationFull(rootTotal)
		payload.SLA["down_time_seconds"] = formatDurationFull(rootDown)
		payload.SLA["up_time_seconds"] = formatDurationFull(rootUp)
		payload.SLA["uptime90"] = fmt.Sprintf("%.3f%%", rootAvail*100)
		payload.SLA["sla_breached"] = (s.Target >= 1.0 && rootDown > 0) || (rootAvail < s.Target)

		idx := -1
		for i, r := range defaultReqs {
			if r.Name == name {
				idx = i
				break
			}
		}

		wrappedPayload := map[string]any{
			"index": idx,
			"payload": map[string]any{
				"probe": payload.Probe,
				"sla":   payload.SLA,
			},
		}

		jsonData, _ := json.Marshal(wrappedPayload)
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		gz.Write(jsonData)
		gz.Close()

		var updateErr error
		if revision > 0 {
			_, updateErr = kv.Update(ctx, name, buf.Bytes(), revision)
		} else {
			_, updateErr = kv.Create(ctx, name, buf.Bytes())
		}

		if updateErr == nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func capSlice[T any](s []T, max int) []T {
	if len(s) > max {
		return s[:max]
	}
	return s
}

func readFromNATS(name string) []byte {

	if nc.Status() != nats.CONNECTED {
		slog.Error("NATS not connected")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kv, err := js.KeyValue(ctx, "BEEP_STATUS")
	if err != nil {
		slog.Error("Failed to access KV bucket", "error", err)
		return nil
	}

	entry, err := kv.Get(ctx, name)
	if err != nil {
		slog.Error("Failed to get entry", "key", name, "error", err)
		return nil
	}

	// Decompress
	reader, err := gzip.NewReader(bytes.NewReader(entry.Value()))
	if err != nil {
		slog.Error("Gzip reader error", "error", err)
		return nil
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		slog.Error("Decompression failed", "error", err)
		return nil
	}

	var wrapped map[string]any
	if err := json.Unmarshal(decompressed, &wrapped); err != nil {
		slog.Error("Unmarshal failed", "error", err)
		return nil
	}

	wrappedData, err := json.Marshal(wrapped)
	if err != nil {
		slog.Error("Marshal failed", "error", err)
		return nil
	}

	return wrappedData

}

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)

	name := r.URL.Query().Get("name")
	history := readFromNATS(name)

	if history == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(history))

}

// -------------------- MAIN --------------------
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	nc, err = nats.Connect(
		serverURL,
		nats.UserJWTAndSeed(jwt, seed),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5*time.Second),
		nats.Timeout(10*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			slog.Warn("Disconnected from NATS", "error", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Info("Reconnected to NATS", "url", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			slog.Error("NATS connection permanently closed")
		}),
	)

	if err != nil {
		slog.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to NATS", "url", serverURL)

	js, err = jetstream.New(nc)

	if err != nil {
		slog.Error("JetStream context error", "error", err)
	}

	kv, err = js.KeyValue(context.Background(), "BEEP_STATUS")
	if err != nil {
		kv, _ = js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
			Bucket:   "BEEP_STATUS",
			MaxBytes: 1024 * 1024 * 50,
		})
	}

	startProbeManager(ctx, &wg)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/sse", Sse)
	mux.HandleFunc("GET /v1/status", StatusHandler)
	mux.HandleFunc("GET /v1/status/history", HistoryHandler)
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", Host, Port),
		Handler: recoveryMiddleware(mux),
	}

	go func() {
		slog.Info("Beep API server running", "url", fmt.Sprintf("http://%s:%s", Host, Port))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			stop()
		}
	}()

	slog.Info("Beep is now active and monitoring services.")

	<-ctx.Done()

	slog.Info("Shutdown signal received. Cleaning up...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	wg.Wait()

	if nc != nil {
		slog.Info("Flushing NATS buffers...")
		if err := nc.Flush(); err != nil {
			slog.Error("NATS flush error", "error", err)
		}
		nc.Close()
		slog.Info("NATS connection closed")
	}

	slog.Info("Shutdown complete. Exiting.")
	os.Exit(0)
}
