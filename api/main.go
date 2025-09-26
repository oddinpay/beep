package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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
	Name     string        `json:"name,omitempty"`
}

type ProbeResult struct {
	Protocol    string `json:"protocol"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Name        string `json:"name,omitempty"`
	Timestamp   string `json:"timestamp"`

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


	resp, err := c.Get(fmt.Sprintf("%s://%s", req.Protocol, req.Host))
	if err != nil {
		return ProbeResult{
			req.Name,
			strings.ToUpper(req.Protocol),
			strings.ToUpper(HealthResponse.Down),
			fmt.Sprintf("%s - %s", req.Host, err.Error()),
			time.Now().Format("15:04:05.000"),
		}
	}

	defer resp.Body.Close()
	return ProbeResult{
		req.Name,
		strings.ToUpper(req.Protocol),
		strings.ToUpper(HealthResponse.Up),
		fmt.Sprintf("%s - %d", req.Host, resp.StatusCode),
		time.Now().Format("15:04:05.000"),

	}

}


func probeTCP(req HttpRequest) ProbeResult {

	
	var HealthResponse = HealthResponse{Down: "down", Up: "up"}

	conn, err := net.DialTimeout("tcp", req.Host, defaultTimeout)
	if err != nil {
		return ProbeResult{
			Protocol:    strings.ToUpper(req.Protocol),
			Status:      strings.ToUpper(HealthResponse.Down),
			Description: err.Error(),
			Name:        req.Name,
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}
	defer conn.Close()

	// Try to write a test byte
	_, err = conn.Write([]byte("ping\n"))
	if err != nil {
		return ProbeResult{
			Protocol:    strings.ToUpper(req.Protocol),
			Status:      strings.ToUpper(HealthResponse.Down),
			Description: "write failed: " + err.Error(),
			Name:        fmt.Sprintf("TCP %s", req.Host),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	buf := make([]byte, 64)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ProbeResult{
			Protocol:    strings.ToUpper(req.Protocol),
			Status:      strings.ToUpper(HealthResponse.Up),
			Description: "no response after connect",
			Name:        fmt.Sprintf("TCP %s", req.Host),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	res := ProbeResult{
		Protocol:    strings.ToUpper(req.Protocol),
		Status:      strings.ToUpper(HealthResponse.Up),
		Description: fmt.Sprintf("response received %s", strings.TrimSpace(string(buf[:n]))),
		Name:        fmt.Sprintf("TCP %s", req.Host),
		Timestamp:   time.Now().Format("15:04:05.000"),
	}
	return res
}


func StatusHandler(w http.ResponseWriter, r *http.Request) {
	reqs := []HttpRequest{
		{Name: "API1", Protocol: "http", Host: "oddinpay.com", Interval: 2 * time.Second},
		{Name: "API2", Protocol: "https",  Host: "github.com", Interval: 20 * time.Second},
		{Name: "API3", Protocol: "tcp",   Host: "localhost:6379"},
	}

	conn, err := sse.Upgrade(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() { _ = conn.Close() }()

	capacity := max(len(reqs) * 10)
	resultChan := make(chan ProbeResult, capacity)


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

		var fn func() ProbeResult
		switch strings.ToLower(strings.TrimSpace(target.Protocol)) {
		case "tcp":
			fn = func() ProbeResult { return probeTCP(target) }
		case "http":
			fn = func() ProbeResult { return probeHTTP(target) }
		case "https":
			fn = func() ProbeResult { return probeHTTP(target) }
		}

		probes = append(probes, probeDef{
			name:     fmt.Sprintf("%s://%s", target.Protocol, target.Host),
			fn:       fn,
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

