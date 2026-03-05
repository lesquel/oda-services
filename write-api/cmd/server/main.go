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

	"github.com/lesquel/oda-write-api/internal/config"
	"github.com/lesquel/oda-write-api/internal/database"
	"github.com/lesquel/oda-write-api/internal/middleware"
	"github.com/lesquel/oda-write-api/internal/seed"

	authhttp "github.com/lesquel/oda-write-api/internal/features/auth/delivery/http"
	authrepo "github.com/lesquel/oda-write-api/internal/features/auth/repository"
	authusecase "github.com/lesquel/oda-write-api/internal/features/auth/usecase"

	poemshttp "github.com/lesquel/oda-write-api/internal/features/poems/delivery/http"
	poemsrepo "github.com/lesquel/oda-write-api/internal/features/poems/repository"
	poemsusecase "github.com/lesquel/oda-write-api/internal/features/poems/usecase"

	adminhttp "github.com/lesquel/oda-write-api/internal/features/admin/delivery/http"
	adminrepo "github.com/lesquel/oda-write-api/internal/features/admin/repository"
	adminusecase "github.com/lesquel/oda-write-api/internal/features/admin/usecase"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.RunMigrations(db); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	if err := seed.SeedAdmin(db, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		slog.Info("admin seed warning", "error", err)
	}
	seed.SeedEmotions(db)

	// ── Repositories ──────────────────────────────────────────────────────
	userRepo := authrepo.NewUserRepository(db)
	tokenRepo := authrepo.NewRefreshTokenRepository(db)
	poemRepo := poemsrepo.NewPoemRepository(db)
	likeRepo := poemsrepo.NewLikeRepository(db)
	emotionRepo := poemsrepo.NewEmotionRepository(db)
	bookmarkRepo := poemsrepo.NewBookmarkRepository(db)
	adminRepoInst := adminrepo.NewAdminRepository(db)

	// ── Use cases ─────────────────────────────────────────────────────────
	authUC := authusecase.NewAuthUseCase(userRepo, tokenRepo, cfg.JWTSecret)
	poemUC := poemsusecase.NewPoemUseCase(poemRepo, likeRepo, emotionRepo, bookmarkRepo)
	adminUC := adminusecase.NewAdminUseCase(adminRepoInst)

	// ── Handlers ──────────────────────────────────────────────────────────
	authH := authhttp.NewAuthHandler(authUC)
	poemH := poemshttp.NewPoemHandler(poemUC)
	adminH := adminhttp.NewAdminHandler(adminUC)

	// ── Router ────────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.InjectUserContext) // always extract user context from gateway headers
	r.Use(middleware.SlogRequestLogger)

	// ── Health check (no auth) ──────────────────────────────────────────
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"write-api"}`))
	})

	// ── Huma API (auto-generates /docs and /openapi.json) ───────────────
	humaConfig := huma.DefaultConfig("ODA Write API", "1.0.0")
	humaConfig.Info.Description = "Handles authentication, poem mutations, bookmarks, likes, and admin operations."
	api := humachi.New(r, humaConfig)

	// ── Huma middleware ─────────────────────────────────────────────────
	internalMW := middleware.HumaInternalAuth(cfg.InternalSecret)
	requireUserMW := middleware.HumaRequireUser(api)
	requireAdminMW := middleware.HumaRequireAdmin(api)

	// ── Register all routes ─────────────────────────────────────────────
	authhttp.RegisterAuthRoutes(api, authH, internalMW, requireUserMW)
	poemshttp.RegisterPoemRoutes(api, poemH, internalMW, requireUserMW)
	adminhttp.RegisterAdminRoutes(api, adminH, internalMW, requireAdminMW)

	// ── Server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("write-api listening", "port", cfg.Port, "docs", "http://localhost:"+cfg.Port+"/docs")
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
	slog.Info("shutting down write-api...")
	_ = srv.Shutdown(ctx)
}
