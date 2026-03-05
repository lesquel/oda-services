package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

type ListBookmarksOutput struct {
	Body domain.PaginatedResponse[domain.AdminBookmark]
}

func (h *AdminHandler) ListBookmarks(ctx context.Context, input *AssociationFilterInput) (*ListBookmarksOutput, error) {
	result, err := h.uc.ListBookmarks(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListBookmarksOutput{Body: *result}, nil
}

func (h *AdminHandler) HardDeleteBookmark(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteBookmark(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
