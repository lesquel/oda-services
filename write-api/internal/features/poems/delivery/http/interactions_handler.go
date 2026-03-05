package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-write-api/internal/middleware"
)

// ── Interaction types ───────────────────────────────────────────────────────

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
	ID string `path:"id" format:"uuid" doc:"Poem UUID"`
}

// ── Interaction handlers ────────────────────────────────────────────────────

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
