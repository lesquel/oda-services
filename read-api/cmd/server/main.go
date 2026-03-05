package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/lesquel/oda-read-api/internal/config"
	"github.com/lesquel/oda-read-api/internal/database"
	"github.com/lesquel/oda-read-api/internal/middleware"

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
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.InjectUserContext)
	r.Use(middleware.SlogRequestLogger)

	// ── Huma API ──────────────────────────────────────────────────────────
	humaConfig := huma.DefaultConfig("ODA Read API", "1.0.0")
	humaConfig.Info.Description = "Read-only API for ODA poetry platform — feed, search, profiles, stats."
	api := humachi.New(r, humaConfig)

	// ── Health check ─────────────────────────────────────────────────────
	huma.Register(api, huma.Operation{
		OperationID: "healthz",
		Summary:     "Health check",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Tags:        []string{"Health"},
	}, func(_ context.Context, _ *struct{}) (*struct {
		Body struct {
			Status string `json:"status"`
		}
	}, error) {
		out := &struct {
			Body struct {
				Status string `json:"status"`
			}
		}{}
		out.Body.Status = "ok"
		return out, nil
	})

	// ── Register read routes ─────────────────────────────────────────────
	internalMW := middleware.HumaInternalAuth(cfg.InternalSecret)
	requireUserMW := middleware.HumaRequireUser(api)
	feedhttp.RegisterReadRoutes(api, h, internalMW, requireUserMW)

	// ── Server ───────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("read-api listening", "port", cfg.Port, "docs", "http://localhost:"+cfg.Port+"/docs")
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
