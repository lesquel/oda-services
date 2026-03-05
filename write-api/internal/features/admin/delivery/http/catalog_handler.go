package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Emotion Catalog types ───────────────────────────────────────────────────

type EmotionCatalogListOutput struct {
	Body []*domain.EmotionCatalog
}

type GetEmotionCatalogOutput struct {
	Body domain.EmotionCatalog
}

type CreateEmotionCatalogInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"50" doc:"Emotion name"`
		Emoji       string `json:"emoji" required:"false" doc:"Emoji representation"`
		Description string `json:"description" required:"false" doc:"Description"`
	}
}

type UpdateEmotionCatalogInput struct {
	ID   string `path:"id" format:"uuid" doc:"Emotion catalog UUID"`
	Body struct {
		Name        string `json:"name,omitempty" minLength:"1" maxLength:"50" required:"false" doc:"Updated name"`
		Emoji       string `json:"emoji,omitempty" required:"false" doc:"Updated emoji"`
		Description string `json:"description,omitempty" required:"false" doc:"Updated description"`
	}
}

// ── Emotion Catalog handlers ────────────────────────────────────────────────

func (h *AdminHandler) ListEmotionCatalog(ctx context.Context, _ *struct{}) (*EmotionCatalogListOutput, error) {
	items, err := h.uc.ListEmotionCatalog()
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &EmotionCatalogListOutput{Body: items}, nil
}

func (h *AdminHandler) GetEmotionCatalog(ctx context.Context, input *GetByIDInput) (*GetEmotionCatalogOutput, error) {
	item, err := h.uc.GetEmotionCatalog(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetEmotionCatalogOutput{Body: *item}, nil
}

func (h *AdminHandler) CreateEmotionCatalog(ctx context.Context, input *CreateEmotionCatalogInput) (*AdminMessageOutput, error) {
	req := &domain.CreateEmotionCatalogRequest{
		Name:        input.Body.Name,
		Emoji:       input.Body.Emoji,
		Description: input.Body.Description,
	}
	if err := h.uc.CreateEmotionCatalog(req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion created"
	return out, nil
}

func (h *AdminHandler) UpdateEmotionCatalog(ctx context.Context, input *UpdateEmotionCatalogInput) (*AdminMessageOutput, error) {
	req := &domain.UpdateEmotionCatalogRequest{
		Name:        input.Body.Name,
		Emoji:       input.Body.Emoji,
		Description: input.Body.Description,
	}
	if err := h.uc.UpdateEmotionCatalog(input.ID, req); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion updated"
	return out, nil
}

func (h *AdminHandler) SoftDeleteEmotionCatalog(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.SoftDeleteEmotionCatalog(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) RestoreEmotionCatalog(ctx context.Context, input *GetByIDInput) (*AdminMessageOutput, error) {
	if err := h.uc.RestoreEmotionCatalog(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion catalog restored"
	return out, nil
}

func (h *AdminHandler) PermanentDeleteEmotionCatalog(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.PermanentDeleteEmotionCatalog(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
