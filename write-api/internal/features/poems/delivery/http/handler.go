package http

import (
	"github.com/lesquel/oda-write-api/internal/features/poems/usecase"
)

// PoemHandler handles poem mutation routes.
type PoemHandler struct{ uc usecase.PoemUseCase }

func NewPoemHandler(uc usecase.PoemUseCase) *PoemHandler { return &PoemHandler{uc: uc} }

// ── Shared types ────────────────────────────────────────────────────────────

type PoemIDInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem UUID"`
}

type PoemMessageOutput struct {
	Body struct {
		Message string `json:"message" doc:"Result message"`
	}
}
