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

	defaultTimeout = 3 * time.Second
)

type HttpRequest struct {
	Host     string `json:"host,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

type ProbeResult struct {
	Protocol    string `json:"protocol"`
	Status      string `json:"status"`
	Description string `json:"description"`
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
	resp, err := c.Get(fmt.Sprintf("%s://%s", "http", "192.168.1.101:5173"))
	if err != nil {
		return ProbeResult{strings.ToUpper(req.Protocol), strings.ToUpper(HealthResponse.Down), err.Error()}
	}
	defer resp.Body.Close()
	return ProbeResult{
		strings.ToUpper(req.Protocol), strings.ToUpper(HealthResponse.Up), fmt.Sprintf("%d", resp.StatusCode)}

}


func StatusHandler(w http.ResponseWriter, r *http.Request) {
    var req HttpRequest = HttpRequest{
        Protocol: "http",
        Host:     "192.168.1.101:5173",
    }

    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(ErrorResponse{
            Status:  "error",
            Message: "Method not allowed",
        })
        return
    }

    // Upgrade the connection to SSE
    conn, err := sse.Upgrade(r.Context(), w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer func() { _ = conn.Close() }()

    // resultChan fan-in for all probes (buffer helps absorb brief bursts)
	resultChan := make(chan ProbeResult, 32)

	// Define the probes you want to run in parallel
	type probeDef struct {
		name string
		fn   func(HttpRequest) ProbeResult
		intravel time.Duration
	}


	probes := []probeDef{
		{name: req.Protocol, fn: probeHTTP, intravel: 1 * time.Second},
		// add more probes here as needed	
	}

	for _, p := range probes {
    go func(p probeDef) {
        ticker := time.NewTicker(p.intravel)
        defer ticker.Stop()

        for {
            select {
            case <-r.Context().Done():
                return
            case <-ticker.C:
                result := p.fn(req)
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

	// Consume and stream results to the SSE connection
	for {
		select {
		case <-r.Context().Done():
			return
		case res := <-resultChan:
			if err := conn.SendData(r.Context(), res); err != nil {
				// Client closed or network error — exit to stop readers via context
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

