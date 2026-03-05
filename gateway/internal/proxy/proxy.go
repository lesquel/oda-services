package proxy

import (
	"bytes"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/lesquel/oda-gateway/internal/middleware"
)

// Config holds proxy configuration.
type Config struct {
	WriteAPIURL    string
	ReadAPIURL     string
	InternalSecret string
	// Client is the HTTP client used for proxying.
	Client *http.Client
}

// Handler wraps the proxy with circuit breakers per upstream.
type Handler struct {
	cfg     Config
	writeCB *CircuitBreaker
	readCB  *CircuitBreaker
}

// New creates a proxy handler with circuit breakers and retry for GETs.
func New(cfg Config) *Handler {
	if cfg.Client == nil {
		cfg.Client = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 25,
				IdleConnTimeout:     90 * time.Second,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ResponseHeaderTimeout: 15 * time.Second,
			},
		}
	}

	return &Handler{
		cfg:     cfg,
		writeCB: NewCircuitBreaker(),
		readCB:  NewCircuitBreaker(),
	}
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetURL := routeRequest(r, h.cfg.WriteAPIURL, h.cfg.ReadAPIURL)
	cb := h.circuitBreakerFor(targetURL)

	if !cb.Allow() {
		slog.Warn("circuit breaker open", "target", targetURL, "uri", r.RequestURI)
		http.Error(w, `{"error":"service temporarily unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	// For idempotent GET requests, retry up to 2 times with exponential backoff
	maxAttempts := 1
	if r.Method == http.MethodGet {
		maxAttempts = 3
	}

	// Buffer the body for non-GET so we can't accidentally re-send mutations
	var bodyBytes []byte
	if r.Body != nil && r.Method != http.MethodGet {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body.Close()
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * 100 * time.Millisecond
			slog.Info("retrying request", "attempt", attempt+1, "backoff_ms", backoff.Milliseconds(), "uri", r.RequestURI)
			time.Sleep(backoff)
		}

		var body io.Reader
		if bodyBytes != nil {
			body = bytes.NewReader(bodyBytes)
		}

		proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL+r.RequestURI, body)
		if err != nil {
			slog.Error("proxy: failed to create request", "error", err, "uri", r.RequestURI)
			http.Error(w, `{"error":"internal proxy error"}`, http.StatusBadGateway)
			return
		}

		copyHeaders(r.Header, proxyReq.Header)
		proxyReq.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)

		if userID := middleware.GetUserID(r.Context()); userID != "" {
			proxyReq.Header.Set("X-User-ID", userID)
			proxyReq.Header.Set("X-User-Role", middleware.GetRole(r.Context()))
		}

		if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
			proxyReq.Header.Set("X-Request-ID", reqID)
		}

		resp, err := h.cfg.Client.Do(proxyReq)
		if err != nil {
			lastErr = err
			cb.RecordFailure()
			slog.Warn("proxy: upstream error", "error", err, "attempt", attempt+1, "target", targetURL)
			continue
		}

		// 5xx from upstream = retryable for GETs
		if resp.StatusCode >= 500 && attempt < maxAttempts-1 && r.Method == http.MethodGet {
			resp.Body.Close()
			cb.RecordFailure()
			slog.Warn("proxy: upstream 5xx", "status", resp.StatusCode, "attempt", attempt+1, "target", targetURL)
			continue
		}

		if resp.StatusCode < 500 {
			cb.RecordSuccess()
		} else {
			cb.RecordFailure()
		}

		// Success or final attempt — send response
		copyHeaders(resp.Header, w.Header())
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
		resp.Body.Close()
		return
	}

	// All attempts exhausted
	slog.Error("proxy: all attempts failed", "error", lastErr, "target", targetURL, "uri", r.RequestURI)
	http.Error(w, `{"error":"service unavailable"}`, http.StatusBadGateway)
}

func (h *Handler) circuitBreakerFor(targetURL string) *CircuitBreaker {
	if targetURL == h.cfg.WriteAPIURL {
		return h.writeCB
	}
	return h.readCB
}

// routeRequest decides which upstream to send the request to.
func routeRequest(r *http.Request, writeURL, readURL string) string {
	if r.Method != http.MethodGet {
		return writeURL
	}

	path := r.URL.Path
	writeGetPrefixes := []string{
		"/api/me",
		"/api/auth/profile",
		"/api/admin/",
	}
	for _, prefix := range writeGetPrefixes {
		if strings.HasPrefix(path, prefix) {
			return writeURL
		}
	}

	return readURL
}

func copyHeaders(src, dst http.Header) {
	for key, values := range src {
		lk := strings.ToLower(key)
		switch {
		case lk == "connection" || lk == "keep-alive" || lk == "proxy-authenticate" ||
			lk == "proxy-authorization" || lk == "te" || lk == "trailers" ||
			lk == "transfer-encoding" || lk == "upgrade":
			continue
		// Strip CORS headers from upstream — the gateway owns CORS.
		case strings.HasPrefix(lk, "access-control-"):
			continue
		}
		for _, v := range values {
			dst.Add(key, v)
		}
	}
}
