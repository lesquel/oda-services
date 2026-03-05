package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

type ListEmotionsOutput struct {
	Body domain.PaginatedResponse[domain.AdminEmotion]
}

func (h *AdminHandler) ListEmotions(ctx context.Context, input *AssociationFilterInput) (*ListEmotionsOutput, error) {
	result, err := h.uc.ListEmotions(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListEmotionsOutput{Body: *result}, nil
}

func (h *AdminHandler) HardDeleteEmotion(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteEmotion(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
