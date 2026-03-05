package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-read-api/internal/middleware"
	"github.com/lesquel/oda-shared/domain"
)

// ── Bookmark types ──────────────────────────────────────────────────────────

type GetBookmarksInput struct {
	Page  int `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit int `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}

type GetBookmarksOutput struct {
	Body []*domain.Poem
}

// ── Bookmark handlers ───────────────────────────────────────────────────────

func (h *ReadHandler) GetUserBookmarks(ctx context.Context, input *GetBookmarksInput) (*GetBookmarksOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, huma.NewError(http.StatusUnauthorized, "authentication required")
	}
	poems, _, err := h.uc.GetUserBookmarks(userID, input.Page, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if poems == nil {
		poems = []*domain.Poem{}
	}
	return &GetBookmarksOutput{Body: poems}, nil
}
