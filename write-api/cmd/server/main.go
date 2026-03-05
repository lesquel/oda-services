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

	"github.com/lesquel/oda-write-api/internal/config"
	"github.com/lesquel/oda-write-api/internal/database"
	"github.com/lesquel/oda-write-api/internal/middleware"
	"github.com/lesquel/oda-write-api/internal/seed"

	// ── features/auth ──────────────────────────────────────────────────────
	authhttp "github.com/lesquel/oda-write-api/internal/features/auth/delivery/http"
	authrepo "github.com/lesquel/oda-write-api/internal/features/auth/repository"
	authusecase "github.com/lesquel/oda-write-api/internal/features/auth/usecase"

	// ── features/poems ─────────────────────────────────────────────────────
	poemshttp "github.com/lesquel/oda-write-api/internal/features/poems/delivery/http"
	poemsrepo "github.com/lesquel/oda-write-api/internal/features/poems/repository"
	poemsusecase "github.com/lesquel/oda-write-api/internal/features/poems/usecase"

	// ── features/admin ─────────────────────────────────────────────────────
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
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.BodyLimit(2 << 20))
	r.Use(chimw.Timeout(30 * time.Second))

	// ── Health check (no auth required) ──────────────────────────────────
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"write-api"}`))
	})

	// All write-api routes require the internal secret header (only gateway can call)
	r.Group(func(r chi.Router) {
		r.Use(middleware.InternalAuth(cfg.InternalSecret))

		r.Route("/api", func(r chi.Router) {
		// ── Me (profile alias) ───────────────────────────────────────────
		r.With(middleware.InjectUserContext, middleware.RequireUser).Get("/me", authH.GetProfile)

		// ── Auth ──────────────────────────────────────────────────────────
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/refresh", authH.Refresh)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Post("/logout", authH.Logout)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Get("/profile", authH.GetProfile)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Put("/profile", authH.UpdateProfile)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Post("/change-password", authH.ChangePassword)
		})

		// ── Users ─────────────────────────────────────────────────────────
		r.With(middleware.InjectUserContext).Get("/users/search", authH.SearchUsers)
		r.With(middleware.InjectUserContext).Get("/users/{username}", authH.GetPublicProfile)

		// ── Poems ─────────────────────────────────────────────────────────
		r.Route("/poems", func(r chi.Router) {
			r.With(middleware.InjectUserContext, middleware.RequireUser).Post("/", poemH.CreatePoem)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Put("/{id}", poemH.UpdatePoem)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Delete("/{id}", poemH.DeletePoem)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Post("/{id}/like", poemH.ToggleLike)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Post("/{id}/bookmark", poemH.ToggleBookmark)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Post("/{id}/emotions", poemH.TagEmotion)
			r.With(middleware.InjectUserContext, middleware.RequireUser).Delete("/{id}/emotions/{emotionID}", poemH.RemoveEmotionTag)
		})

		// ── Admin ─────────────────────────────────────────────────────────
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.InjectUserContext, middleware.RequireAdmin)

			r.Get("/stats", adminH.GetStats)

			r.Route("/users", func(r chi.Router) {
				r.Get("/", adminH.ListUsers)
				r.Post("/", adminH.CreateUser)
				r.Get("/{id}", adminH.GetUser)
				r.Put("/{id}", adminH.UpdateUser)
				r.Patch("/{id}/role", adminH.ChangeUserRole)
				r.Delete("/{id}", adminH.HardDeleteUser)
			})

			r.Route("/poems", func(r chi.Router) {
				r.Get("/", adminH.ListPoems)
				r.Get("/{id}", adminH.GetPoem)
				r.Put("/{id}", adminH.UpdatePoem)
				r.Patch("/{id}/status", adminH.ChangePoemStatus)
				r.Delete("/{id}", adminH.HardDeletePoem)
			})

			r.Route("/likes", func(r chi.Router) {
				r.Get("/", adminH.ListLikes)
				r.Delete("/{id}", adminH.HardDeleteLike)
			})

			r.Route("/bookmarks", func(r chi.Router) {
				r.Get("/", adminH.ListBookmarks)
				r.Delete("/{id}", adminH.HardDeleteBookmark)
			})

			r.Route("/emotions", func(r chi.Router) {
				r.Get("/", adminH.ListEmotions)
				r.Delete("/{id}", adminH.HardDeleteEmotion)
			})

			r.Route("/emotion-catalog", func(r chi.Router) {
				r.Get("/", adminH.ListEmotionCatalog)
				r.Post("/", adminH.CreateEmotionCatalog)
				r.Put("/{id}", adminH.UpdateEmotionCatalog)
				r.Delete("/{id}", adminH.DeleteEmotionCatalog)
			})
		})
	})
	}) // end r.Group (InternalAuth)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("write-api listening", "port", cfg.Port)
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
