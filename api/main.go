package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.jetify.com/sse"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/fetch"

	"github.com/oklog/ulid/v2"
)

const (
	Host = "127.0.0.1"
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

	defaultTimeout = 10 * time.Second
	minutes90d     = 90 * 24 * 60
)

// ----------- DB / CACHE CONNECTIONS -----------

var (
	// redisClient      valkey.Client
	// fs               *bigcache.BigCache
	probeManagerOnce sync.Once
	probeUpdates     = make(chan map[string]StatusPayload, 100)
)

// -------------------- GLOBAL SLA MAP --------------------

var slaTrackers = struct {
	sync.Mutex
	m map[string]*SlidingSLA
}{m: make(map[string]*SlidingSLA)}

var defaultReqs = func() []HttpRequest {
	raw := []HttpRequest{
		{Name: "HTTPS", Protocol: "https", Host: "oddinpay.com", Interval: 30 * time.Second},
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

// func initRedis() {
// 	var err error
// 	redisClient, err = valkey.NewClient(valkey.ClientOption{
// 		InitAddress: []string{"localhost:6379"},
// 		Username:    "",
// 		Password:    "",
// 	})
// 	if err != nil {
// 		log.Printf("⚠️ Redis unavailable, continuing without cache: %v", err)
// 		redisClient = nil
// 	}
// }

// func initBigcache() {
// 	ctx := context.Background()

// 	cache, err := bigcache.New(ctx, bigcache.DefaultConfig(5*time.Second))
// 	if err != nil {
// 		log.Fatalf("failed to init BigCache: %v", err)
// 		fs = nil
// 	}
// 	fs = cache

// }

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

// -------------------- PROBES --------------------

func probeHTTP(re HttpRequest) ProbeResult {

	c := fetch.NewClient()

	r, err := fetch.NewRequest(context.Background(), http.MethodGet, fmt.Sprintf("%s://%s", re.Protocol, re.Host), nil)

	resp, err := c.Do(r, nil)

	if err != nil {
		return ProbeResult{
			Id:          ulid.Make().String(),
			Name:        re.Name,
			Protocol:    strings.ToUpper(re.Protocol),
			Description: fmt.Sprintf("%s - %s", re.Host, err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
			State:       []string{},
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < StatusOK || resp.StatusCode >= StatusBadRequest {
		return ProbeResult{
			Id:          ulid.Make().String(),
			Name:        re.Name,
			Protocol:    strings.ToUpper(re.Protocol),
			Description: fmt.Sprintf("%s - %d", re.Host, resp.StatusCode),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
			State:       []string{},
		}
	}
	return ProbeResult{
		Id:          ulid.Make().String(),
		Name:        re.Name,
		Protocol:    strings.ToUpper(re.Protocol),
		Description: fmt.Sprintf("%s - %d", re.Host, resp.StatusCode),
		Timestamp:   time.Now().Format("15:04:05.000"),
		Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
		State:       []string{hr.Up},
	}
}

func probeTCP(req HttpRequest) ProbeResult {
	conn, err := net.DialTimeout("tcp", req.Host, defaultTimeout)

	if err != nil {
		return ProbeResult{
			Id:          ulid.Make().String(),
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			Description: err.Error(),
			Timestamp:   time.Now().Format("15:04:05.000"),
			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
			State:       []string{hr.Down, "up", "up", "up"},
		}
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping\n"))
	if err != nil {
		return ProbeResult{
			Id:          ulid.Make().String(),
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
			Id:          ulid.Make().String(),
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Up},
			Description: "no response after connect",
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	return ProbeResult{
		Id:          ulid.Make().String(),
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		State:       []string{hr.Up},
		Description: fmt.Sprintf("response received %s", strings.TrimSpace(string(buf[:n]))),
		Timestamp:   time.Now().Format("15:04:05.000"),
	}
}

func probeDNS(req HttpRequest) ProbeResult {

	if net.ParseIP(req.Host) != nil {
		return ProbeResult{
			Id:          ulid.Make().String(),
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Down},
			Description: "Input is already an IP, DNS lookup skipped",
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	addrs, err := net.LookupHost(req.Host)
	if err != nil {
		return ProbeResult{
			Id:          ulid.Make().String(),
			Name:        req.Name,
			Protocol:    strings.ToUpper(req.Protocol),
			State:       []string{hr.Down},
			Description: fmt.Sprintf("DNS error: %s", err.Error()),
			Timestamp:   time.Now().Format("15:04:05.000"),
		}
	}

	return ProbeResult{
		Id:          ulid.Make().String(),
		Name:        req.Name,
		Protocol:    strings.ToUpper(req.Protocol),
		Description: fmt.Sprintf("resolved %v", addrs),
		Timestamp:   time.Now().Format("15:04:05.000"),
		Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
		State:       []string{hr.Up, "warn", "up", "up"},
	}
}

// func probeUDP(req HttpRequest) ProbeResult {

// 	raddr, err := net.ResolveUDPAddr("udp", req.Host)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    strings.ToUpper(req.Protocol),
// 			Description: "Error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{hr.Down, "down", "up", "down"},
// 		}
// 	}

// 	conn, err := net.DialUDP("udp", nil, raddr)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    strings.ToUpper(req.Protocol),
// 			Description: "dial error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{hr.Down, "down", "up", "down"},
// 		}
// 	}
// 	defer conn.Close()
// 	_ = conn.SetDeadline(time.Now().Add(defaultTimeout))

// 	_, err = conn.Write([]byte("ping\n"))
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    strings.ToUpper(req.Protocol),
// 			Description: "write error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{hr.Down, "down", "up", "down"},
// 		}
// 	}

// 	// Try read (optional)
// 	buf := make([]byte, 64)
// 	n, _, err := conn.ReadFromUDP(buf)
// 	if err != nil {
// 		// No reply → still count as UP
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    strings.ToUpper(req.Protocol),
// 			Description: "write ok (no reply)",
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{"warn", hr.Up, "up", "down"},
// 		}
// 	}

// 	return ProbeResult{
// 		Id:          ulid.Make().String(),
// 		Name:        req.Name,
// 		Protocol:    strings.ToUpper(req.Protocol),
// 		Description: fmt.Sprintf("response received %s", strings.TrimSpace(string(buf[:n]))),
// 		Timestamp:   time.Now().Format("15:04:05.000"),
// 		Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 		State:       []string{hr.Up, "warn"},
// 	}
// }

// func ProbeICMP(req HttpRequest) ProbeResult {

// 	pinger, err := probing.NewPinger(req.Host)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "ICMP",
// 			Description: "Pinger error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{hr.Down},
// 		}
// 	}
// 	pinger.Count = 1
// 	pinger.Timeout = defaultTimeout
// 	pinger.SetPrivileged(true)

// 	err = pinger.Run()
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "ICMP",
// 			Description: "Run error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{hr.Down},
// 		}
// 	}
// 	stats := pinger.Statistics()

// 	if stats.PacketsRecv == 0 {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "ICMP",
// 			State:       []string{hr.Up},
// 			Description: fmt.Sprintf("0/%d packets received", stats.PacketsSent),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	return ProbeResult{
// 		Id:          ulid.Make().String(),
// 		Name:        req.Name,
// 		Protocol:    "ICMP",
// 		State:       []string{hr.Up},
// 		Description: fmt.Sprintf("%d/%d packets received, avg rtt %.2fms", stats.PacketsRecv, stats.PacketsSent, float64(stats.AvgRtt.Microseconds())/1000.0),
// 		Timestamp:   time.Now().Format("15:04:05.000"),
// 	}
// }

// func probeSMTP(req HttpRequest) ProbeResult {

// 	// extract host and port (default to 25 if missing)
// 	hostOnly, port, err := net.SplitHostPort(req.Host)
// 	if err != nil {
// 		hostOnly = req.Host
// 		port = "25"
// 	}

// 	// Use Cloudflare DNS (1.1.1.1) for resolution via a custom resolver
// 	resolver := &net.Resolver{
// 		PreferGo: true,
// 		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
// 			d := net.Dialer{Timeout: defaultTimeout}
// 			return d.DialContext(ctx, "udp", "1.1.1.1:53")
// 		},
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
// 	defer cancel()

// 	ips, err := resolver.LookupIPAddr(ctx, hostOnly)
// 	if err != nil || len(ips) == 0 {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "SMTP",
// 			State:       []string{hr.Down},
// 			Description: "DNS lookup (Cloudflare) failed: " + fmt.Sprintf("%v", err),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// pick first resolved IP and dial that IP:port
// 	targetAddr := net.JoinHostPort(ips[0].IP.String(), port)
// 	dialer := net.Dialer{Timeout: defaultTimeout}
// 	conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "SMTP",
// 			State:       []string{hr.Down},
// 			Description: "Dial failed to resolved IP: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	c, err := smtp.NewClient(conn, hostOnly)
// 	if err != nil {
// 		conn.Close()
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "SMTP",
// 			State:       []string{hr.Down},
// 			Description: "NewClient failed: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}
// 	defer c.Close()

// 	hostname, _ := os.Hostname()

// 	if err := c.Hello(hostname); err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    "SMTP",
// 			State:       []string{hr.Down},
// 			Description: "EHLO failed: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// STARTTLS upgrade
// 	if ok, _ := c.Extension("STARTTLS"); ok {
// 		tlsConfig := &tls.Config{ServerName: hostOnly}
// 		if err = c.StartTLS(tlsConfig); err != nil {
// 			return ProbeResult{
// 				Id:          ulid.Make().String(),
// 				Name:        req.Name,
// 				Protocol:    "SMTP",
// 				State:       []string{hr.Down},
// 				Description: "STARTTLS failed: " + err.Error(),
// 				Timestamp:   time.Now().Format("15:04:05.000"),
// 			}
// 		}
// 	}

// 	desc := "Connected to " + targetAddr + " (resolved via Cloudflare) without authentication"

// 	// AUTH if username/password are set
// 	if strings.TrimSpace(req.Username) != "" && strings.TrimSpace(req.Password) != "" {
// 		if ok, _ := c.Extension("AUTH"); !ok {
// 			return ProbeResult{
// 				Id:          ulid.Make().String(),
// 				Name:        req.Name,
// 				Protocol:    "SMTP",
// 				State:       []string{hr.Down},
// 				Description: "Server does not support AUTH",
// 				Timestamp:   time.Now().Format("15:04:05.000"),
// 			}
// 		}

// 		auth := smtp.PlainAuth("", req.Username, req.Password, hostOnly)
// 		if err := c.Auth(auth); err != nil {
// 			return ProbeResult{
// 				Id:          ulid.Make().String(),
// 				Name:        req.Name,
// 				Protocol:    "SMTP",
// 				State:       []string{hr.Down},
// 				Description: "AUTH failed: " + err.Error(),
// 				Timestamp:   time.Now().Format("15:04:05.000"),
// 			}
// 		}
// 		desc = fmt.Sprintf("Authenticated to %s successfully (resolved via Cloudflare)", req.Host)
// 	}

// 	return ProbeResult{
// 		Id:          ulid.Make().String(),
// 		Name:        req.Name,
// 		Protocol:    "SMTP",
// 		State:       []string{hr.Up},
// 		Description: desc,
// 		Timestamp:   time.Now().Format("15:04:05.000"),
// 	}
// }

// func ProbeRedis(req HttpRequest) ProbeResult {

// 	opt := valkey.ClientOption{
// 		InitAddress: []string{req.Host},
// 	}

// 	if strings.TrimSpace(req.Username) != "" {
// 		opt.Username = req.Username
// 	}
// 	if strings.TrimSpace(req.Password) != "" {
// 		opt.Password = req.Password
// 	}

// 	client, err := valkey.NewClient(opt)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			Description: "Client init error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 			Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 			State:       []string{hr.Down, "up", "up", "up"},
// 		}
// 	}
// 	defer client.Close()

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// 1) PING
// 	replyPing, err := client.Do(ctx, client.B().Ping().Build()).ToString()
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "PING error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 2) SET key
// 	replySet, err := client.Do(ctx, client.B().Set().Key("c9289d4f-8ff8-412c-bf3a-9d59a9776979").Value("OK").Build()).ToString()
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "SET error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 3) GET key
// 	val, err := client.Do(ctx, client.B().Get().Key("c9289d4f-8ff8-412c-bf3a-9d59a9776979").Build()).ToString()
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "GET error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 4) DEL key (cleanup, ignore errors)
// 	_, _ = client.Do(ctx, client.B().Del().Key("c9289d4f-8ff8-412c-bf3a-9d59a9776979").Build()).AsInt64()

// 	// 5) Return probe result
// 	desc := fmt.Sprintf("ping:%s, set:%s, get:%s", replyPing, replySet, val)

// 	return ProbeResult{
// 		Id:          ulid.Make().String(),
// 		Name:        req.Name,
// 		Protocol:    req.Protocol,
// 		Description: desc,
// 		Timestamp:   time.Now().Format("15:04:05.000"),
// 		Date:        []string{time.Now().Format("02/01/2006"), "29/09/2025", "26/09/2025", "25/09/2025"},
// 		State:       []string{hr.Up, "warn", "up", "up"},
// 	}

// }

// func ProbePostgres(req HttpRequest) ProbeResult {

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 	defer cancel()

// 	conn, err := pgx.Connect(ctx, req.Protocol+"://"+req.Host)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Connect error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}
// 	defer conn.Close(ctx)

// 	// 1) Ensure users table exists
// 	_, err = conn.Exec(ctx, `
// 		create table if not exists users (
// 			id serial primary key,
// 			username text not null unique,
// 			age int
// 		)
// 	`)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Create table error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 2) Insert user
// 	_, err = conn.Exec(ctx, "insert into users (username, age) values ($1, $2) on conflict (username) do nothing", "jack", 30)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Insert error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 3) Select user
// 	var age int
// 	err = conn.QueryRow(ctx, "select age from users where username=$1", "jack").Scan(&age)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Select error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 4) Update user
// 	_, err = conn.Exec(ctx, "update users set age=$1 where username=$2", 31, "jack")
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Update error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 5) Select again
// 	err = conn.QueryRow(ctx, "select age from users where username=$1", "jack").Scan(&age)
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Re-select error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 6) Delete user
// 	_, err = conn.Exec(ctx, "delete from users where username=$1", "jack")
// 	if err != nil {
// 		return ProbeResult{
// 			Id:          ulid.Make().String(),
// 			Name:        req.Name,
// 			Protocol:    req.Protocol,
// 			State:       []string{hr.Down},
// 			Description: "Delete error: " + err.Error(),
// 			Timestamp:   time.Now().Format("15:04:05.000"),
// 		}
// 	}

// 	// 7) Verify delete
// 	// err = conn.QueryRow(ctx, "select age from users where username=$1", "jack").Scan(&age)
// 	// if err != nil {
// 	// 	fmt.Println("After delete → no such user (as expected)")
// 	// }

// 	desc := fmt.Sprintf("Connection is working - Querying: age:%d:ok", age)
// 	return ProbeResult{
// 		Id:          ulid.Make().String(),
// 		Name:        req.Name,
// 		Protocol:    req.Protocol,
// 		State:       []string{hr.Up},
// 		Description: desc,
// 		Timestamp:   time.Now().Format("15:04:05.000"),
// 	}
// }

// -------------------- 90-DAY SLA --------------------

func NewSlidingSLA(target float64) *SlidingSLA {
	now := time.Now().Truncate(time.Minute)
	return &SlidingSLA{
		Target:        target,
		buckets:       make([]bucket, minutes90d),
		currentMinute: now,
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

// 🔹 Tick now uses real elapsed wall time
func (s *SlidingSLA) Tick(isDown bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.rotateTo(now)

	elapsed := max(int64(now.Sub(s.lastUpdate).Seconds()), 0)
	s.lastUpdate = now

	s.buckets[s.idx].totalSec += elapsed
	if isDown {
		s.buckets[s.idx].downSec += elapsed
	}
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
			"id":                 ulid.Make().String(),
			"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
			"uptime90":           "100.000%",
			"up_time_seconds":    formatDurationFull(0),
			"down_time_seconds":  formatDurationFull(0),
			"total_time_seconds": formatDurationFull(0),
			"sla_breached":       false,
		}
	}

	availability := 1.0 - (float64(down) / float64(total))
	breached := (s.Target >= 1.0 && down > 0)
	up := total - down

	return map[string]any{
		"id":                 ulid.Make().String(),
		"sla_target":         fmt.Sprintf("%.3f%%", s.Target*100),
		"uptime90":           fmt.Sprintf("%.3f%%", availability*100),
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

// Starts global probe manager once
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
			// case "udp":
			// 	probeFn = probeUDP
			// case "smtp":
			// 	probeFn = probeSMTP
			// case "redis":
			// 	probeFn = ProbeRedis
			// case "postgres":
			// 	probeFn = ProbePostgres
			// case "icmp":
			// 	probeFn = ProbeICMP
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
					tracker.Tick(isDown)

					payload := StatusPayload{
						Probe: res,
						SLA:   tracker.Snapshot(),
					}

					// Write to Redis / BigCache
					// data, _ := json.Marshal(payload)
					// key := fmt.Sprintf("probe:%s:%s", req.Name, time.Now().Format("2006-01-02"))
					// if redisClient != nil {
					// 	_ = redisClient.Do(
					// 		ctx,
					// 		redisClient.B().Set().
					// 			Key(key).
					// 			Value(string(data)).
					// 			Ex(90*24*time.Hour).
					// 			Build(),
					// 	)
					// }
					// if fs != nil {
					// 	_ = fs.Set(key, data)
					// }

					// Broadcast update
					select {
					case probeUpdates <- map[string]StatusPayload{req.Name: payload}:
					default: // avoid blocking if no listener
					}

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

	ctx := r.Context()

	// Send cached data immediately
	// for i, req := range defaultReqs {
	// 	key := fmt.Sprintf("probe:%s:%s", req.Name, time.Now().Format("2006-01-02"))
	// 	if payload, err := loadPayload(ctx, key); err == nil {
	// 		out := map[string]any{
	// 			"index":   i,
	// 			"payload": *payload,
	// 		}
	// 		_ = conn.SendData(ctx, out)
	// 	}
	// }

	// Stream live updates from global channel
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-probeUpdates:
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
					if !errors.Is(err, context.Canceled) {
						log.Printf("⚠️ SSE send error [%s]: %v", name, err)
					}
					return
				}
			}
		}
	}
}

// func loadPayload(ctx context.Context, key string) (*StatusPayload, error) {
// 	if fs != nil {
// 		if cached, err := fs.Get(key); err == nil {
// 			var p StatusPayload
// 			if jsonErr := json.Unmarshal(cached, &p); jsonErr == nil {
// 				return &p, nil
// 			}
// 		}
// 	}
// 	// if redisClient != nil {
// 	// 	val, err := redisClient.Do(ctx, redisClient.B().Get().Key(key).Build()).ToString()
// 	// 	if err == nil && val != "" {
// 	// 		var p StatusPayload
// 	// 		if jsonErr := json.Unmarshal([]byte(val), &p); jsonErr == nil {
// 	// 			return &p, nil
// 	// 		}
// 	// 	}
// 	// }
// 	return nil, errors.New("cache miss")
// }

// -------------------- STATE REQUEST HANDLER --------------------

func RestRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != MethodGet {
		http.Error(w, "Method not allowed", StatusMethodNotAllowed)
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

// -------------------- MAIN --------------------

func main() {

	// initRedis()
	// initBigcache()

	// if redisClient != nil {
	// 	defer redisClient.Close()
	// }

	// if fs != nil {
	// 	defer fs.Close()
	// }

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/sse", StatusHandler)
	mux.HandleFunc("/v1/status", RestRequestHandler)
	mux.HandleFunc("/v1/sla/reset", ResetHandler)

	handler := recoveryMiddleware(mux)

	fmt.Printf("Beep API server running at http://%s:%s\n", Host, Port)
	// if err := http.ListenAndServe(fmt.Sprintf("%s:%s", Host, Port), handler); err != nil {
	// 	log.Fatal(err)
	// }
	//
	workers.Serve(handler)

}
