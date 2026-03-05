package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Bookmark types ──────────────────────────────────────────────────────────

type ListBookmarksOutput struct {
	Body domain.PaginatedResponse[domain.AdminBookmark]
}

type GetBookmarkOutput struct {
	Body domain.AdminBookmark
}

type CreateBookmarkInput struct {
	Body struct {
		UserID string `json:"user_id" format:"uuid" doc:"User UUID"`
		PoemID string `json:"poem_id" format:"uuid" doc:"Poem UUID"`
	}
}

// ── Bookmark handlers ───────────────────────────────────────────────────────

func (h *AdminHandler) ListBookmarks(ctx context.Context, input *AssociationFilterInput) (*ListBookmarksOutput, error) {
	result, err := h.uc.ListBookmarks(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListBookmarksOutput{Body: *result}, nil
}

func (h *AdminHandler) GetBookmark(ctx context.Context, input *GetByIDInput) (*GetBookmarkOutput, error) {
	bookmark, err := h.uc.GetBookmark(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetBookmarkOutput{Body: *bookmark}, nil
}

func (h *AdminHandler) CreateBookmark(ctx context.Context, input *CreateBookmarkInput) (*AdminMessageOutput, error) {
	req := &domain.CreateBookmarkRequest{
		UserID: input.Body.UserID,
		PoemID: input.Body.PoemID,
	}
	if err := h.uc.CreateBookmark(req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "bookmark created"
	return out, nil
}

func (h *AdminHandler) SoftDeleteBookmark(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.SoftDeleteBookmark(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) RestoreBookmark(ctx context.Context, input *GetByIDInput) (*AdminMessageOutput, error) {
	if err := h.uc.RestoreBookmark(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "bookmark restored"
	return out, nil
}

func (h *AdminHandler) PermanentDeleteBookmark(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.PermanentDeleteBookmark(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
