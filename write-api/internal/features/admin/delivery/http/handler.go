package http

import (
	"github.com/lesquel/oda-write-api/internal/features/admin/usecase"
)

// AdminHandler handles admin routes.
type AdminHandler struct{ uc usecase.AdminUseCase }

func NewAdminHandler(uc usecase.AdminUseCase) *AdminHandler { return &AdminHandler{uc: uc} }

// ── Shared input / output types ─────────────────────────────────────────────

type PaginatedInput struct {
	Page  int    `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	Q     string `query:"q" required:"false" doc:"Search query"`
}

type GetByIDInput struct {
	ID string `path:"id" format:"uuid" doc:"Resource UUID"`
}

type AssociationFilterInput struct {
	Page   int    `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	PoemID string `query:"poem_id" required:"false" doc:"Filter by poem UUID"`
	UserID string `query:"user_id" required:"false" doc:"Filter by user UUID"`
}

type AdminMessageOutput struct {
	Body struct {
		Message string `json:"message" doc:"Result message"`
	}
}
