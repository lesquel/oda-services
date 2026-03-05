package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Poem types ──────────────────────────────────────────────────────────────

type GetFeedInput struct {
	Page  int `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit int `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type GetFeedOutput struct {
	Body []*domain.Poem
}

type SearchPoemsInput struct {
	Q     string `query:"q" required:"false" doc:"Search query"`
	Page  int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type SearchPoemsOutput struct {
	Body []*domain.Poem
}

type GetPoemInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem ID"`
}
type GetPoemOutput struct {
	Body *domain.Poem
}

type GetPoemStatsInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem ID"`
}
type GetPoemStatsOutput struct {
	Body struct {
		LikesCount   int   `json:"likes_count" doc:"Total likes"`
		ViewsCount   int   `json:"views_count" doc:"Total views"`
		EmotionCount int64 `json:"emotion_count" doc:"Total emotion tags"`
	}
}

type GetEmotionDistInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem ID"`
}
type GetEmotionDistOutput struct {
	Body map[string]int
}

type GetUserPoemsInput struct {
	UserID string `path:"userID" format:"uuid" doc:"User ID"`
	Status string `query:"status" required:"false" enum:"published,draft" default:"published" doc:"Filter by poem status"`
	Page   int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type GetUserPoemsOutput struct {
	Body []*domain.Poem
}

// ── Poem handlers ───────────────────────────────────────────────────────────

func (h *ReadHandler) GetFeed(_ context.Context, input *GetFeedInput) (*GetFeedOutput, error) {
	poems, _, err := h.uc.GetFeed(input.Page, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if poems == nil {
		poems = []*domain.Poem{}
	}
	return &GetFeedOutput{Body: poems}, nil
}

func (h *ReadHandler) SearchPoems(_ context.Context, input *SearchPoemsInput) (*SearchPoemsOutput, error) {
	poems, _, err := h.uc.SearchPoems(input.Q, input.Page, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if poems == nil {
		poems = []*domain.Poem{}
	}
	return &SearchPoemsOutput{Body: poems}, nil
}

func (h *ReadHandler) GetPoem(_ context.Context, input *GetPoemInput) (*GetPoemOutput, error) {
	poem, err := h.uc.GetPoem(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "poem not found")
	}
	return &GetPoemOutput{Body: poem}, nil
}

func (h *ReadHandler) GetPoemStats(_ context.Context, input *GetPoemStatsInput) (*GetPoemStatsOutput, error) {
	stats, err := h.uc.GetPoemStats(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	out := &GetPoemStatsOutput{}
	if v, ok := stats["likes_count"].(int); ok {
		out.Body.LikesCount = v
	}
	if v, ok := stats["views_count"].(int); ok {
		out.Body.ViewsCount = v
	}
	if v, ok := stats["emotion_count"].(int64); ok {
		out.Body.EmotionCount = v
	}
	return out, nil
}

func (h *ReadHandler) GetEmotionDistribution(_ context.Context, input *GetEmotionDistInput) (*GetEmotionDistOutput, error) {
	dist, err := h.uc.GetEmotionDistribution(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetEmotionDistOutput{Body: dist}, nil
}

func (h *ReadHandler) GetUserPoems(_ context.Context, input *GetUserPoemsInput) (*GetUserPoemsOutput, error) {
	poems, _, err := h.uc.GetUserPoems(input.UserID, input.Status, input.Page, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if poems == nil {
		poems = []*domain.Poem{}
	}
	return &GetUserPoemsOutput{Body: poems}, nil
}
