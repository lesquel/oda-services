package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Emotion catalog types ───────────────────────────────────────────────────

type GetEmotionCatalogOutput struct {
	Body []*domain.EmotionCatalog
}

// ── Emotion catalog handlers ────────────────────────────────────────────────

func (h *ReadHandler) GetEmotionCatalog(_ context.Context, _ *struct{}) (*GetEmotionCatalogOutput, error) {
	items, err := h.uc.GetEmotionCatalog()
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if items == nil {
		items = []*domain.EmotionCatalog{}
	}
	return &GetEmotionCatalogOutput{Body: items}, nil
}
