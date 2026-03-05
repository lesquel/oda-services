package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Like types ──────────────────────────────────────────────────────────────

type ListLikesOutput struct {
	Body domain.PaginatedResponse[domain.AdminLike]
}

type GetLikeOutput struct {
	Body domain.AdminLike
}

type CreateLikeInput struct {
	Body struct {
		UserID string `json:"user_id" format:"uuid" doc:"User UUID"`
		PoemID string `json:"poem_id" format:"uuid" doc:"Poem UUID"`
	}
}

// ── Like handlers ───────────────────────────────────────────────────────────

func (h *AdminHandler) ListLikes(ctx context.Context, input *AssociationFilterInput) (*ListLikesOutput, error) {
	result, err := h.uc.ListLikes(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListLikesOutput{Body: *result}, nil
}

func (h *AdminHandler) GetLike(ctx context.Context, input *GetByIDInput) (*GetLikeOutput, error) {
	like, err := h.uc.GetLike(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetLikeOutput{Body: *like}, nil
}

func (h *AdminHandler) CreateLike(ctx context.Context, input *CreateLikeInput) (*AdminMessageOutput, error) {
	req := &domain.CreateLikeRequest{
		UserID: input.Body.UserID,
		PoemID: input.Body.PoemID,
	}
	if err := h.uc.CreateLike(req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "like created"
	return out, nil
}

func (h *AdminHandler) SoftDeleteLike(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.SoftDeleteLike(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) RestoreLike(ctx context.Context, input *GetByIDInput) (*AdminMessageOutput, error) {
	if err := h.uc.RestoreLike(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "like restored"
	return out, nil
}

func (h *AdminHandler) PermanentDeleteLike(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.PermanentDeleteLike(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
