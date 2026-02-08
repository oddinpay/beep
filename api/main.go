package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"go.jetify.com/sse"
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
	minutes90d     = 90 * 24 * 60
)

// ----------- DB / CACHE CONNECTIONS -----------

var (
	jwt              = os.Getenv("NATS_JWT")
	seed             = os.Getenv("NATS_SEED")
	serverURL        = os.Getenv("NATS_URL")
	probeManagerOnce sync.Once

	monitorStartTime = time.Now().UTC().Truncate(24 * time.Hour)
	nc               = func() *nats.Conn {
		c, err := nats.Connect(serverURL, nats.UserJWTAndSeed(jwt, seed))
		if err != nil {
			slog.Error("Failed to connect to NATS server", "error", err)
			os.Exit(1)
		}
		return c
	}()
)

// -------------------- GLOBAL SLA MAP --------------------

var slaTrackers = struct {
	sync.Mutex
	m map[string]*SlidingSLA
}{m: make(map[string]*SlidingSLA)}

var defaultReqs = func() []HttpRequest {
	raw := []HttpRequest{
		{Name: "DNS", Protocol: "dns", Host: "www.oddinpay.com", Interval: 10 * time.Second},
	}

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

var hr = HealthResponse{Down: "down", Up: "up"}

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
	todayUTC := time.Now().UTC().Truncate(24 * time.Hour)
	daysSinceStart := int(todayUTC.Sub(monitorStartTime).Hours()/24) + 1

	count := min(daysSinceStart, 90)

	dates := make([]string, count)
	for i := range count {
		dates[i] = todayUTC.AddDate(0, 0, -i).Format("02/01/2006")
	}
	return dates
}

// -------------------- BROADCAST HUB --------------------

type Hub struct {
	sync.RWMutex
	clients map[chan map[string]StatusPayload]struct{}
}

var globalHub = &Hub{
	clients: make(map[chan map[string]StatusPayload]struct{}),
}

func (h *Hub) Broadcast(update map[string]StatusPayload) {
	h.RLock()
	defer h.RUnlock()
	for clientChan := range h.clients {
		select {
		case clientChan <- update:
		default:
		}
	}
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
	r.Header.Set("User-Agent", "beep_01kgwc0fggeze9075f1tk43bdf/1.0")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return ProbeResult{
			Name:        re.Name,
			Protocol:    strings.ToUpper(re.Protocol),
			Description: fmt.Sprintf("%s - %s", re.Host, err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{},
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < StatusOK || resp.StatusCode >= StatusBadRequest {
		return ProbeResult{
			Name:        re.Name,
			Protocol:    strings.ToUpper(re.Protocol),
			Description: fmt.Sprintf("%s - %d", re.Host, resp.StatusCode),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{},
		}
	}
	return ProbeResult{
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
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: err.Error(),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        getRecentDates(),
			State:       []string{hr.Down, "up", "up", "up"},
		}
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping\n"))
	if err != nil {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Down},
			Description: "write failed: " + err.Error(),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	buf := make([]byte, 64)
	_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Up},
			Description: "no response after connect",
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	return ProbeResult{
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		State:       []string{hr.Up},
		Description: fmt.Sprintf("response received %s", strings.TrimSpace(string(buf[:n]))),
		Timestamp:   time.Now().Format("15:04:05.000"),
	}
}

func probeDNS(req HttpRequest) ProbeResult {

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if net.ParseIP(req.Host) != nil {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Down},
			Description: "Input is already an IP, DNS lookup skipped",
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	addrs, err := net.DefaultResolver.LookupHost(ctx, req.Host)
	if err != nil {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Down},
			Description: fmt.Sprintf("DNS error: %s", err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	return ProbeResult{
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		Description: fmt.Sprintf("resolved %v", addrs),
		Timestamp:   time.Now().Format("15:04:05.000"),
		Date:        getRecentDates(),
		State:       []string{hr.Up},
	}
}

// -------------------- 90-DAY SLA --------------------

func NewSlidingSLA(target float64) *SlidingSLA {
	now := time.Now()
	return &SlidingSLA{
		Target:        target,
		buckets:       make([]bucket, minutes90d),
		currentMinute: now.Truncate(time.Minute),
		lastUpdate:    now,
	}
}

func (s *SlidingSLA) rotateTo(now time.Time) {
	minNow := now.Truncate(time.Minute)
	if !minNow.After(s.currentMinute) {
		return
	}
	steps := int(minNow.Sub(s.currentMinute) / time.Minute)
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

// -------------------- SSE HANDLER --------------------

func startProbeManager() {
	probeManagerOnce.Do(func() {
		log.Println("🚀 Starting global probe manager...")

		for _, target := range defaultReqs {
			t := target

			interval := t.Interval
			if interval <= 0 {
				interval = 1 * time.Second
			}

			var probeFn func(HttpRequest) ProbeResult
			switch strings.ToLower(strings.TrimSpace(t.Protocol)) {
			case "tcp":
				probeFn = probeTCP
			case "http", "https":
				probeFn = probeHTTP
			case "dns":
				probeFn = probeDNS
			default:
				log.Printf("⚠️ Unsupported protocol: %s", t.Protocol)
				continue
			}

			go func(req HttpRequest, fn func(HttpRequest) ProbeResult, iv time.Duration) {

				ticker := time.NewTicker(iv)
				defer ticker.Stop()

				for range ticker.C {
					// ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

					_, cancel := context.WithTimeout(context.Background(), defaultTimeout)

					res := fn(req)

					slaTrackers.Lock()
					tracker := slaTrackers.m[req.Name]
					if tracker == nil {
						tracker = NewSlidingSLA(1.0)
						slaTrackers.m[req.Name] = tracker
					}
					slaTrackers.Unlock()

					isDown := len(res.State) == 0 || strings.ToLower(res.State[0]) != "up"
					tracker.Tick(isDown, interval)

					payload := StatusPayload{
						Probe: res,
						SLA:   tracker.Snapshot(),
					}

					publishToNATS(req.Name, payload)

					// Broadcast update
					globalHub.Broadcast(map[string]StatusPayload{req.Name: payload})

					cancel()
				}
			}(t, probeFn, interval)
		}
	})
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != MethodPost && r.Method != MethodGet {
		http.Error(w, "Method not allowed", StatusMethodNotAllowed)
		return
	}

	// SSE headers
	w.Header().Set(HeaderAllowOrigin, "*")
	w.Header().Set(HeaderCacheControl, "no-cache")
	w.Header().Set(HeaderConnection, "keep-alive")
	w.Header().Set(HeaderContentType, ContentTypeEventStream)

	// Start probe manager if not already running
	startProbeManager()

	conn, err := sse.Upgrade(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), StatusInternalServerError)
		return
	}
	defer conn.Close()

	clientChan := make(chan map[string]StatusPayload, 50)

	globalHub.Lock()
	globalHub.clients[clientChan] = struct{}{}
	globalHub.Unlock()

	defer func() {
		globalHub.Lock()
		delete(globalHub.clients, clientChan)
		globalHub.Unlock()
	}()

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-clientChan:
			for name, payload := range update {
				idx := -1
				for i, r := range defaultReqs {
					if r.Name == name {
						idx = i
						break
					}
				}
				out := map[string]any{
					"index":   idx,
					"payload": payload,
				}
				if err := conn.SendData(ctx, out); err != nil {
					return
				}
			}
		}
	}
}

// -------------------- STATE REQUEST HANDLER --------------------

func RestRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != MethodGet {
		http.Error(w, "Method not allowed", StatusMethodNotAllowed)
		return
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)

	var hasMonitors bool
	var miniMonitors bool
	if len(defaultReqs) == 0 {
		hasMonitors = false
	} else {
		hasMonitors = true
	}

	if len(defaultReqs) > 2 {
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

func publishToNATS(name string, payload StatusPayload) {
	if nc.Status() != nats.CONNECTED {
		slog.Error("NATS not connected")
	}

	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("JetStream context error", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	streamName := "STATUS"
	subject := fmt.Sprintf("STATUS.%s", name)

	s, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  jetstream.FileStorage,
		MaxBytes: 1024 * 1024 * 50,
	})

	data, _ := json.Marshal(payload)

	ack, err := js.Publish(ctx, subject, data)
	if err != nil {
		slog.Error("Publish failed", "error", err)
	}

	fmt.Printf("Appended to %s | Seq: %d\n", ack.Stream, ack.Sequence)

	c, _ := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	msgs, err := c.FetchNoWait(2)
	if err != nil {
		slog.Error("Fetch failed", "error", err)
	}

	for msg := range msgs.Messages() {
		fmt.Printf("Received a JetStream message: %s\n", string(msg.Data()))
	}
}

// -------------------- MAIN --------------------
func main() {

	startProbeManager()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/sse", StatusHandler)
	mux.HandleFunc("/v1/status", RestRequestHandler)
	mux.HandleFunc("/v1/sla/reset", ResetHandler)

	handler := recoveryMiddleware(mux)

	fmt.Printf("Beep API server running at http://%s:%s\n", Host, Port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", Host, Port), handler); err != nil {
		slog.Error("Server failed to start", "error", err)
	}

}
