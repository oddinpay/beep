package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.jetify.com/sse"
)

const (
	Host = "0.0.0.0"
	Port = "8976"

	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
	MethodPatch  = "PATCH"

	ContentTypeJSON = "application/json"

	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusMethodNotAllowed    = 405
	StatusMultipleChoices     = 300

	defaultTimeout = 60 * time.Second
)

type HttpRequest struct {
	Host     string        `json:"host,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
}

type ProbeResult struct {
	Protocol    string `json:"protocol"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Timestamp string `json:"timestamp"`

}

type HealthResponse struct {
	Down string `json:"down"`
	Up   string `json:"up"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Recovery middleware
func recoveryMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v\n", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{
					Status:  "error",
					Message: "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func probeHTTP(req HttpRequest) ProbeResult {

	var HealthResponse = HealthResponse{Down: "down", Up: "up"}
	c := &http.Client{Timeout: defaultTimeout}

	protocol := strings.ToLower(strings.TrimSpace(req.Protocol))

	host := strings.TrimSpace(req.Host)
	host = strings.TrimPrefix(host, ".")
	host = strings.TrimPrefix(host, "*.")

	parts := strings.FieldsFunc(host, func(r rune) bool { return r == '.' })
	base := ""

	
	if len(parts) >= 2 {
		base = parts[len(parts)-2]
	} else if len(parts) == 1 {
		base = parts[0]
	}
	base = strings.TrimSpace(base)
	if base != "" {
		base = strings.ToUpper(base[:1]) + strings.ToLower(base[1:])
	}

	if protocol != "http" && protocol != "https" {
		return ProbeResult{
			strings.ToUpper(req.Protocol),
			strings.ToUpper(HealthResponse.Down),
			fmt.Sprintf("%s - %d protocol not allowed", base, http.StatusMethodNotAllowed),
			time.Now().Format("16:04:05.000"),
		}
	}

	resp, err := c.Get(fmt.Sprintf("%s://%s", protocol, req.Host))
	if err != nil {
		return ProbeResult{
			strings.ToUpper(req.Protocol),
			strings.ToUpper(HealthResponse.Down),
			fmt.Sprintf("%s - %s", base, err.Error()),
			time.Now().Format("15:04:05.000"),
		}
	}

	defer resp.Body.Close()
	return ProbeResult{
		strings.ToUpper(req.Protocol),
		strings.ToUpper(HealthResponse.Up),
		fmt.Sprintf("%s - %d", base, resp.StatusCode),
		time.Now().Format("15:04:05.000"),

	}

}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	reqs := []HttpRequest{
		{Protocol: "https", Host: "oddinpay.com", Interval: 2 * time.Second},
		{Protocol: "http", Host: "github.com", Interval: 20 * time.Second},
	}

	conn, err := sse.Upgrade(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() { _ = conn.Close() }()

	resultChan := make(chan ProbeResult, 32)

	type probeDef struct {
		name     string
		fn       func() ProbeResult
		intravel time.Duration
	}

	probes := make([]probeDef, 0, len(reqs))
	for _, t := range reqs {
		target := t
		interval := target.Interval
		if interval <= 0 {
			interval = 1 * time.Second
		}
		probes = append(probes, probeDef{
			name:     fmt.Sprintf("%s://%s", target.Protocol, target.Host),
			fn:       func() ProbeResult { return probeHTTP(target) },
			intravel: interval,
		})
	}

	for _, p := range probes {
		go func(p probeDef) {
			
			initial := p.fn()
			select {
			case resultChan <- initial:
			case <-r.Context().Done():
				return
			default:
			}

			ticker := time.NewTicker(p.intravel)
			defer ticker.Stop()

			for {
				select {
				case <-r.Context().Done():
					return
				case <-ticker.C:
					result := p.fn()
					select {
					case resultChan <- result:
					case <-r.Context().Done():
						return
					default:
					}
				}
			}
		}(p)
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case res := <-resultChan:
			if err := conn.SendData(r.Context(), res); err != nil {
				return
			}
		}
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/status", StatusHandler)

	// Wrap the mux with recovery
	handler := recoveryMiddleware(mux)

	fmt.Printf("API server running at http://%s:%s\n", Host, Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", Host, Port), handler); err != nil {
		return
	}

}

