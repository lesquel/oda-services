package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Emotion types ───────────────────────────────────────────────────────────

type ListEmotionsOutput struct {
	Body domain.PaginatedResponse[domain.AdminEmotion]
}

type GetEmotionOutput struct {
	Body domain.AdminEmotion
}

type CreateEmotionInput struct {
	Body struct {
		UserID    string `json:"user_id" format:"uuid" doc:"User UUID"`
		PoemID    string `json:"poem_id" format:"uuid" doc:"Poem UUID"`
		EmotionID string `json:"emotion_id" format:"uuid" doc:"Emotion catalog UUID"`
	}
}

// ── Emotion handlers ────────────────────────────────────────────────────────

func (h *AdminHandler) ListEmotions(ctx context.Context, input *AssociationFilterInput) (*ListEmotionsOutput, error) {
	result, err := h.uc.ListEmotions(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListEmotionsOutput{Body: *result}, nil
}

func (h *AdminHandler) GetEmotion(ctx context.Context, input *GetByIDInput) (*GetEmotionOutput, error) {
	emotion, err := h.uc.GetEmotion(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetEmotionOutput{Body: *emotion}, nil
}

func (h *AdminHandler) CreateEmotion(ctx context.Context, input *CreateEmotionInput) (*AdminMessageOutput, error) {
	req := &domain.CreateEmotionTagRequest{
		UserID:    input.Body.UserID,
		PoemID:    input.Body.PoemID,
		EmotionID: input.Body.EmotionID,
	}
	if err := h.uc.CreateEmotion(req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion tag created"
	return out, nil
}

func (h *AdminHandler) SoftDeleteEmotion(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.SoftDeleteEmotion(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) RestoreEmotion(ctx context.Context, input *GetByIDInput) (*AdminMessageOutput, error) {
	if err := h.uc.RestoreEmotion(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion tag restored"
	return out, nil
}

func (h *AdminHandler) PermanentDeleteEmotion(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.PermanentDeleteEmotion(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
