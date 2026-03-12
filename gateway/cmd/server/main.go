package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/lesquel/oda-gateway/internal/config"
	"github.com/lesquel/oda-gateway/internal/middleware"
	"github.com/lesquel/oda-gateway/internal/proxy"
)

func main() {
	// ── Structured logging ──────────────────────────────────────────────
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.Load()
	docsClient := &http.Client{Timeout: 15 * time.Second}

	proxyHandler := proxy.New(proxy.Config{
		WriteAPIURL:    cfg.WriteAPIURL,
		ReadAPIURL:     cfg.ReadAPIURL,
		InternalSecret: cfg.InternalSecret,
	})

	r := chi.NewRouter()

	// ── Global middleware ────────────────────────────────────────────────
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(middleware.SlogRequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ── Security headers ────────────────────────────────────────────────
	r.Use(middleware.SecurityHeaders)

	// ── Health check (not proxied) ──────────────────────────────────────
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"gateway"}`))
	})

	// ── API docs index ──────────────────────────────────────────────────
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html><head><title>ODA API Docs</title></head><body style="font-family:system-ui;max-width:600px;margin:80px auto">
<h1>ODA API Documentation</h1>
<ul>
<li><a href="/docs/write">Write API (mutations)</a></li>
<li><a href="/docs/read">Read API (queries)</a></li>
</ul>
</body></html>`))
	})

	// ── Docs proxied through gateway (same origin :8080) ───────────────
	r.Get("/docs/write", proxyDocsPage(docsClient, cfg.WriteAPIURL, "/docs/write/openapi.yaml"))
	r.Get("/docs/read", proxyDocsPage(docsClient, cfg.ReadAPIURL, "/docs/read/openapi.yaml"))
	r.Get("/docs/write/openapi.yaml", proxyDocsSpec(docsClient, cfg.WriteAPIURL))
	r.Get("/docs/read/openapi.yaml", proxyDocsSpec(docsClient, cfg.ReadAPIURL))

	// ── API routes ──────────────────────────────────────────────────────
	r.Route("/api", func(r chi.Router) {
		// Public routes — optional JWT (user context injected if present)
		r.Group(func(r chi.Router) {
			r.Use(middleware.OptionalJWTAuth(cfg.JWTSecret))
			r.Handle("/poems", proxyHandler)
			r.Handle("/poems/*", proxyHandler)
			r.Handle("/users", proxyHandler)
			r.Handle("/users/*", proxyHandler)
			r.Handle("/emotions", proxyHandler)
			r.Handle("/emotions/*", proxyHandler)
		})

		// Auth routes — no auth required (register, login, refresh)
		r.Handle("/auth/register", proxyHandler)
		r.Handle("/auth/login", proxyHandler)
		r.Handle("/auth/refresh", proxyHandler)

		// Authenticated routes — JWT required
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireJWTAuth(cfg.JWTSecret))
			r.Handle("/me", proxyHandler)
			r.Handle("/auth/*", proxyHandler)
			r.Handle("/bookmarks", proxyHandler)
			r.Handle("/bookmarks/*", proxyHandler)
		})

		// Admin routes — admin JWT required
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdminAuth(cfg.JWTSecret))
			r.Handle("/admin", proxyHandler)
			r.Handle("/admin/*", proxyHandler)
		})
	})

	// ── Server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("gateway listening", "port", cfg.Port, "write_api", cfg.WriteAPIURL, "read_api", cfg.ReadAPIURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	slog.Info("shutting down gateway...")
	_ = srv.Shutdown(ctx)
}

func proxyDocsPage(client *http.Client, upstreamBaseURL, openapiPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, strings.TrimRight(upstreamBaseURL, "/")+"/docs", nil)
		if err != nil {
			http.Error(w, `{"error":"failed to build docs request"}`, http.StatusBadGateway)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, `{"error":"failed to fetch docs"}`, http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, `{"error":"failed to read docs response"}`, http.StatusBadGateway)
			return
		}

		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(strings.ToLower(contentType), "text/html") {
			body = []byte(strings.ReplaceAll(string(body), "/openapi.yaml", openapiPath))
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(body)
	}
}

func proxyDocsSpec(client *http.Client, upstreamBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, strings.TrimRight(upstreamBaseURL, "/")+"/openapi.yaml", nil)
		if err != nil {
			http.Error(w, `{"error":"failed to build openapi request"}`, http.StatusBadGateway)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, `{"error":"failed to fetch openapi"}`, http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		if contentType := resp.Header.Get("Content-Type"); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}
