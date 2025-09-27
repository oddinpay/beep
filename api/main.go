package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

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


	// HTTP status codes
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

	defaultTimeout = 5 * time.Second
	minutes90d    = 90 * 24 * 60
	secondsPerMin = 60
)


// ----------- CORS whitelist (edit as needed) -----------

var allowedOrigins = []string{
	"https://app1.local",
}


// -------------------- GLOBAL SLA MAP --------------------

var slaTrackers = struct {
	sync.Mutex
	m map[string]*SlidingSLA
}{m: make(map[string]*SlidingSLA)}


// -------------------- MODELS --------------------

type HttpRequest struct {
	Host     string        `json:"host,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
	Name     string        `json:"name,omitempty"`
}

type HealthResponse struct {
	Down string `json:"down"`
	Up   string `json:"up"`
}

type ProbeResult struct {
	Name        string `json:"name,omitempty"`
	Protocol    string `json:"protocol"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

type StatusPayload struct {
	Probe     ProbeResult    `json:"probe"`
	SLA       map[string]any `json:"sla"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}


type bucket struct{ totalSec, downSec int32 }


type SlidingSLA struct {
	Target        float64
	buckets       []bucket
	idx           int
	currentMinute time.Time
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
					Status:  "error",
					Message: "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// -------------------- CORS --------------------

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowOrigin := ""
		for _, ao := range allowedOrigins {
			if ao == origin {
				allowOrigin = ao
				break
			}
		}
		if allowOrigin == "" {
			http.Error(w, "CORS origin not allowed", StatusUnauthorized)
			return
		}

		w.Header().Set(HeaderAllowOrigin, allowOrigin)
		w.Header().Set(HeaderAllowMethods, "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set(HeaderAllowHeaders, "Content-Type, Authorization")

		if r.Method == MethodOptions {
			w.WriteHeader(StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// -------------------- PROBES --------------------

func probeHTTP(req HttpRequest) ProbeResult {
	var hr = HealthResponse{Down: "down", Up: "up"}
	c := &http.Client{Timeout: defaultTimeout}

	resp, err := c.Get(fmt.Sprintf("%s://%s", req.Protocol, req.Host))
	if err != nil {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Status:      strings.ToUpper(hr.Down),
			Description: fmt.Sprintf("%s - %s", req.Host, err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}
	defer resp.Body.Close()

	return ProbeResult{
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		Status:      strings.ToUpper(hr.Up),
		Description: fmt.Sprintf("%s - %d", req.Host, resp.StatusCode),
		Timestamp:   time.Now().Format("15:04:05.000"),
	}
}

func probeTCP(req HttpRequest) ProbeResult {
	var hr = HealthResponse{Down: "down", Up: "up"}

	conn, err := net.DialTimeout("tcp", req.Host, defaultTimeout)


	if err != nil {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Status:      strings.ToUpper(hr.Down),
			Description: err.Error(),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping\n"))
	if err != nil {
		return ProbeResult{
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Status:      strings.ToUpper(hr.Down),
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
			Status:      strings.ToUpper(hr.Up),
			Description: "no response after connect",
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	return ProbeResult{
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		Status:      strings.ToUpper(hr.Up),
		Description: fmt.Sprintf("response received %s", strings.TrimSpace(string(buf[:n]))),
		Timestamp:   time.Now().Format("15:04:05.000"),
	}
}

// -------------------- 90-DAY SLA --------------------


func NewSlidingSLA(target float64) *SlidingSLA {
	now := time.Now()
	return &SlidingSLA{
		Target:        target,
		buckets:       make([]bucket, minutes90d),
		currentMinute: now.Truncate(time.Minute),
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
	for i := 0; i < steps; i++ {
		s.idx++
		if s.idx >= len(s.buckets) {
			s.idx = 0
		}
		s.buckets[s.idx] = bucket{}
	}
	s.currentMinute = minNow
}


func (s *SlidingSLA) Tick(isDown bool) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rotateTo(now)
	if s.buckets[s.idx].totalSec < secondsPerMin {
		s.buckets[s.idx].totalSec++
	}
	if isDown && s.buckets[s.idx].downSec < secondsPerMin {
		s.buckets[s.idx].downSec++
	}
}


func (s *SlidingSLA) Snapshot() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	var total, down int64
	for _, b := range s.buckets {
		total += int64(b.totalSec)
		down += int64(b.downSec)
	}

	if total <= 0 {
		return map[string]any{
			"sla_target":             "100.000%",
			"availability_percent":   "100.000%",
			"up_time_seconds":        0,
			"down_time_seconds":      0,
			"total_time_seconds":     0,
			"sla_breached":           false,
		}
	}

	availability := 1.0 - (float64(down) / float64(total))
	breached := (s.Target >= 1.0 && down > 0)
	up := total - down

	return map[string]any{
		"sla_target":           fmt.Sprintf("%.3f%%", s.Target*100),
		"availability_percent": fmt.Sprintf("%.3f%%", availability*100),
		"up_time_seconds":      up,
		"down_time_seconds":    down,
		"total_time_seconds":   total,
		"sla_breached":         breached,
	}
}


func (s *SlidingSLA) Reset() {
	s.mu.Lock()
	for i := range s.buckets {
		s.buckets[i] = bucket{}
	}
	s.idx = 0
	s.currentMinute = time.Now().Truncate(time.Minute)
	s.mu.Unlock()
}



// -------------------- SSE HANDLER --------------------

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderCacheControl, "no-cache")
	w.Header().Set(HeaderConnection, "keep-alive")
	w.Header().Set(HeaderContentType, ContentTypeEventStream)

	reqs := []HttpRequest{
		{Name: "", 	   Protocol: "ht",  Host: "oddinpay.com"},
		{Name: "API2", Protocol: "hts", Host: "github.com", Interval: 10 * time.Second},
		{Name: "API3", Protocol: "t",   Host: "localhost:6379"},
	}

	conn, err := sse.Upgrade(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), StatusInternalServerError)
		return
	}

	slaTrackers.Lock()
	for _, t := range reqs {
		if _, ok := slaTrackers.m[t.Name]; !ok {
			slaTrackers.m[t.Name] = NewSlidingSLA(1.0)
		}
	}
	slaTrackers.Unlock()

	for _, t := range reqs {
		target := t
		interval := target.Interval
		if interval <= 0 {
			interval = 1 * time.Second
		}

		var fn func(HttpRequest) ProbeResult
		switch strings.ToLower(strings.TrimSpace(target.Protocol)) {
		case "tcp":
			fn = probeTCP
		case "http", "https":
			fn = probeHTTP
		default:
			continue
		}

		go func(req HttpRequest, pfn func(HttpRequest) ProbeResult, iv time.Duration) {
			ticker := time.NewTicker(iv)
			defer ticker.Stop()
			for {
				select {
				case <-r.Context().Done():
					return
				case <-ticker.C:
					res := pfn(req)

					slaTrackers.Lock()
					tracker := slaTrackers.m[req.Name]
					slaTrackers.Unlock()
					isDown := strings.ToUpper(res.Status) != "UP"
					tracker.Tick(isDown)

					payload := StatusPayload{
						Probe:     res,
						SLA:       tracker.Snapshot(),
					}
					_ = conn.SendData(r.Context(), payload)
				}
			}
		}(target, fn, interval)
	}

	<-r.Context().Done()
}

// -------------------- RESET HANDLER --------------------

func ResetHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"sla_reset": true,
		"probe":     name,
	})
}

// -------------------- MAIN --------------------

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/uptime", StatusHandler)
	mux.HandleFunc("/sla/reset", ResetHandler)

	handler := recoveryMiddleware(corsMiddleware(mux))

	fmt.Printf("API server running at http://%s:%s\n", Host, Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", Host, Port), handler); err != nil {
		log.Fatal(err)
	}
}

