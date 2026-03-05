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

	"github.com/lesquel/oda-read-api/internal/config"
	"github.com/lesquel/oda-read-api/internal/database"
	"github.com/lesquel/oda-read-api/internal/middleware"

	// ── features/feed ──────────────────────────────────────────────────────
	feedhttp "github.com/lesquel/oda-read-api/internal/features/feed/delivery/http"
	feedrepo "github.com/lesquel/oda-read-api/internal/features/feed/repository"
	feeduc "github.com/lesquel/oda-read-api/internal/features/feed/usecase"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// ── Repository → UseCase → Handler ────────────────────────────────────
	repo := feedrepo.NewReadRepository(db)
	uc := feeduc.NewReadUseCase(repo)
	h := feedhttp.NewReadHandler(uc)

	// ── Router ────────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.BodyLimit(1 << 20))
	r.Use(chimw.Timeout(30 * time.Second))

	// ── Health check (no auth required) ──────────────────────────────────
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"read-api"}`))
	})

	// All read-api routes require the internal secret header (only gateway can call)
	r.Group(func(r chi.Router) {
		r.Use(middleware.InternalAuth(cfg.InternalSecret))

		r.Route("/api", func(r chi.Router) {
			r.Route("/poems", func(r chi.Router) {
				r.With(middleware.InjectUserContext).Get("/feed", h.GetFeed)
				r.With(middleware.InjectUserContext).Get("/search", h.SearchPoems)
				r.With(middleware.InjectUserContext).Get("/{id}", h.GetPoem)
				r.With(middleware.InjectUserContext).Get("/{id}/stats", h.GetPoemStats)
				r.With(middleware.InjectUserContext).Get("/{id}/emotions/distribution", h.GetEmotionDistribution)
			})

			r.Route("/users", func(r chi.Router) {
				r.With(middleware.InjectUserContext).Get("/search", h.SearchUsers)
				r.With(middleware.InjectUserContext).Get("/{username}", h.GetPublicProfile)
				r.With(middleware.InjectUserContext).Get("/{userID}/poems", h.GetUserPoems)
			})

			r.With(middleware.InjectUserContext).Get("/bookmarks", h.GetUserBookmarks)
			r.Get("/emotions", h.GetEmotionCatalog)
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("read-api listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	slog.Info("shutting down read-api...")
	_ = srv.Shutdown(ctx)
}
