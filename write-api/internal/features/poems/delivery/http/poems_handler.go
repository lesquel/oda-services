package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-write-api/internal/middleware"
)

// ── Poem types ──────────────────────────────────────────────────────────────

type CreatePoemInput struct {
	Body struct {
		Title   string `json:"title" minLength:"1" maxLength:"200" doc:"Poem title"`
		Content string `json:"content" minLength:"1" doc:"Poem content (plain text or markdown)"`
		Status  string `json:"status,omitempty" enum:"published,draft" default:"published" required:"false" doc:"Poem status"`
	}
}
type CreatePoemOutput struct {
	Body domain.Poem
}

type UpdatePoemInput struct {
	ID   string `path:"id" format:"uuid" doc:"Poem UUID"`
	Body struct {
		Title   string `json:"title,omitempty" minLength:"1" maxLength:"200" required:"false" doc:"Updated title"`
		Content string `json:"content,omitempty" minLength:"1" required:"false" doc:"Updated content"`
		Status  string `json:"status,omitempty" enum:"published,draft" required:"false" doc:"Updated status"`
	}
}
type UpdatePoemOutput struct {
	Body domain.Poem
}

// ── Poem handlers ───────────────────────────────────────────────────────────

func (h *PoemHandler) CreatePoem(ctx context.Context, input *CreatePoemInput) (*CreatePoemOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	req := &domain.CreatePoemRequest{
		Title:   input.Body.Title,
		Content: input.Body.Content,
		Status:  input.Body.Status,
	}
	poem, err := h.uc.CreatePoem(userID, req)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &CreatePoemOutput{Body: *poem}, nil
}

func (h *PoemHandler) UpdatePoem(ctx context.Context, input *UpdatePoemInput) (*UpdatePoemOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	req := &domain.UpdatePoemRequest{
		Title:   input.Body.Title,
		Content: input.Body.Content,
		Status:  input.Body.Status,
	}
	poem, err := h.uc.UpdatePoem(input.ID, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "unauthorized to update this poem" {
			code = http.StatusForbidden
		}
		return nil, huma.NewError(code, err.Error())
	}
	return &UpdatePoemOutput{Body: *poem}, nil
}

func (h *PoemHandler) DeletePoem(ctx context.Context, input *PoemIDInput) (*struct{}, error) {
	userID, _ := middleware.GetUserID(ctx)
	if err := h.uc.DeletePoem(input.ID, userID); err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "unauthorized to delete this poem" {
			code = http.StatusForbidden
		}
		return nil, huma.NewError(code, err.Error())
	}
	return nil, nil
}
