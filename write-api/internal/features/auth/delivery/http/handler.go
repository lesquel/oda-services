package http

import (
	"github.com/lesquel/oda-write-api/internal/features/auth/usecase"
)

// AuthHandler handles auth and user mutation routes.
type AuthHandler struct{ uc usecase.AuthUseCase }

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler { return &AuthHandler{uc: uc} }

// ── Shared types ────────────────────────────────────────────────────────────

type MessageOutput struct {
	Body struct {
		Message string `json:"message" doc:"Response message"`
	}
}
