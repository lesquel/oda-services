package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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
<li><a href="` + cfg.WriteAPIURL + `/docs">Write API (mutations)</a></li>
<li><a href="` + cfg.ReadAPIURL + `/docs">Read API (queries)</a></li>
</ul>
</body></html>`))
	})

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
