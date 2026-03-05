package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

type ListLikesOutput struct {
	Body domain.PaginatedResponse[domain.AdminLike]
}

func (h *AdminHandler) ListLikes(ctx context.Context, input *AssociationFilterInput) (*ListLikesOutput, error) {
	result, err := h.uc.ListLikes(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListLikesOutput{Body: *result}, nil
}

func (h *AdminHandler) HardDeleteLike(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteLike(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
