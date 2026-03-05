package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-write-api/internal/features/poems/usecase"
	"github.com/lesquel/oda-write-api/internal/middleware"
)

// PoemHandler handles poem mutation routes.
type PoemHandler struct{ uc usecase.PoemUseCase }

func NewPoemHandler(uc usecase.PoemUseCase) *PoemHandler { return &PoemHandler{uc: uc} }

// ── Input / Output types ────────────────────────────────────────────────────

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

type PoemIDInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem UUID"`
}

type ToggleLikeOutput struct {
	Body struct {
		IsLiked bool `json:"is_liked" doc:"Whether the poem is now liked"`
	}
}

type ToggleBookmarkOutput struct {
	Body struct {
		Bookmarked bool `json:"bookmarked" doc:"Whether the poem is now bookmarked"`
	}
}

type TagEmotionInput struct {
	ID   string `path:"id" format:"uuid" doc:"Poem UUID"`
	Body struct {
		EmotionID string `json:"emotion_id" format:"uuid" doc:"Emotion catalog UUID"`
	}
}

type RemoveEmotionInput struct {
	ID        string `path:"id" format:"uuid" doc:"Poem UUID"`
	EmotionID string `path:"emotionID" format:"uuid" doc:"Emotion tag UUID"`
}

type PoemMessageOutput struct {
	Body struct {
		Message string `json:"message" doc:"Result message"`
	}
}

// ── Handlers ────────────────────────────────────────────────────────────────

func (h *PoemHandler) CreatePoem(ctx context.Context, input *CreatePoemInput) (*CreatePoemOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	req := &domain.CreatePoemRequest{
		Title:   input.Body.Title,
		Content: input.Body.Content,
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

func (h *PoemHandler) ToggleLike(ctx context.Context, input *PoemIDInput) (*ToggleLikeOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	liked, err := h.uc.ToggleLike(input.ID, userID)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &ToggleLikeOutput{}
	out.Body.IsLiked = liked
	return out, nil
}

func (h *PoemHandler) ToggleBookmark(ctx context.Context, input *PoemIDInput) (*ToggleBookmarkOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	bookmarked, err := h.uc.ToggleBookmark(input.ID, userID)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &ToggleBookmarkOutput{}
	out.Body.Bookmarked = bookmarked
	return out, nil
}

func (h *PoemHandler) TagEmotion(ctx context.Context, input *TagEmotionInput) (*PoemMessageOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	if err := h.uc.TagEmotion(input.ID, userID, input.Body.EmotionID); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &PoemMessageOutput{}
	out.Body.Message = "emotion tagged"
	return out, nil
}

func (h *PoemHandler) RemoveEmotionTag(ctx context.Context, input *RemoveEmotionInput) (*struct{}, error) {
	userID, _ := middleware.GetUserID(ctx)
	if err := h.uc.RemoveEmotionTag(input.ID, userID); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

// RegisterRoutes registers all poem-related Huma operations.
func RegisterPoemRoutes(api huma.API, h *PoemHandler, internalMW, requireUserMW func(huma.Context, func(huma.Context))) {
	authMW := huma.Middlewares{internalMW, requireUserMW}

	huma.Register(api, huma.Operation{
		OperationID:   "create-poem",
		Summary:       "Create a new poem",
		Method:        http.MethodPost,
		Path:          "/api/poems",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Poems"},
		Middlewares:    authMW,
	}, h.CreatePoem)

	huma.Register(api, huma.Operation{
		OperationID: "update-poem",
		Summary:     "Update an existing poem",
		Method:      http.MethodPut,
		Path:        "/api/poems/{id}",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.UpdatePoem)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-poem",
		Summary:       "Delete a poem",
		Method:        http.MethodDelete,
		Path:          "/api/poems/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Poems"},
		Middlewares:    authMW,
	}, h.DeletePoem)

	huma.Register(api, huma.Operation{
		OperationID: "toggle-like",
		Summary:     "Like or unlike a poem",
		Method:      http.MethodPost,
		Path:        "/api/poems/{id}/like",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.ToggleLike)

	huma.Register(api, huma.Operation{
		OperationID: "toggle-bookmark",
		Summary:     "Bookmark or unbookmark a poem",
		Method:      http.MethodPost,
		Path:        "/api/poems/{id}/bookmark",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.ToggleBookmark)

	huma.Register(api, huma.Operation{
		OperationID: "tag-emotion",
		Summary:     "Tag an emotion on a poem",
		Method:      http.MethodPost,
		Path:        "/api/poems/{id}/emotions",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.TagEmotion)

	huma.Register(api, huma.Operation{
		OperationID:   "remove-emotion-tag",
		Summary:       "Remove your emotion tag from a poem",
		Method:        http.MethodDelete,
		Path:          "/api/poems/{id}/emotions/{emotionID}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Poems"},
		Middlewares:    authMW,
	}, h.RemoveEmotionTag)
}
