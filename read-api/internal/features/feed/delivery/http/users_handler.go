package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── User types ──────────────────────────────────────────────────────────────

type GetPublicProfileInput struct {
	Username string `path:"username" doc:"Username to look up"`
}
type GetPublicProfileOutput struct {
	Body *domain.User
}

type SearchUsersInput struct {
	Q      string `query:"q" required:"false" doc:"Search query"`
	Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Max results"`
	Offset int    `query:"offset" minimum:"0" default:"0" doc:"Offset for pagination"`
}
type SearchUsersOutput struct {
	Body []*domain.User
}

type GetUserStatsInput struct {
	UserID string `path:"userID" format:"uuid" doc:"User ID"`
}
type GetUserStatsOutput struct {
	Body struct {
		PoemCount      int64          `json:"poem_count" doc:"Total poems"`
		PublishedCount int64          `json:"published_count" doc:"Published poems"`
		DraftCount     int64          `json:"draft_count" doc:"Draft poems"`
		TotalLikes     int64          `json:"total_likes" doc:"Total likes received"`
		TotalViews     int64          `json:"total_views" doc:"Total views received"`
		TotalBookmarks int64          `json:"total_bookmarks" doc:"Total bookmarks received"`
		EmotionDist    map[string]int `json:"emotion_distribution" doc:"Emotion distribution across poems"`
	}
}

// ── User handlers ───────────────────────────────────────────────────────────

func (h *ReadHandler) GetPublicProfile(_ context.Context, input *GetPublicProfileInput) (*GetPublicProfileOutput, error) {
	user, err := h.uc.GetPublicProfile(input.Username)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "user not found")
	}
	return &GetPublicProfileOutput{Body: user}, nil
}

func (h *ReadHandler) SearchUsers(_ context.Context, input *SearchUsersInput) (*SearchUsersOutput, error) {
	users, err := h.uc.SearchUsers(input.Q, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if users == nil {
		users = []*domain.User{}
	}
	return &SearchUsersOutput{Body: users}, nil
}

func (h *ReadHandler) GetUserStats(_ context.Context, input *GetUserStatsInput) (*GetUserStatsOutput, error) {
	stats, err := h.uc.GetUserStats(input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &GetUserStatsOutput{}
	if v, ok := stats["poem_count"].(int64); ok {
		out.Body.PoemCount = v
	}
	if v, ok := stats["published_count"].(int64); ok {
		out.Body.PublishedCount = v
	}
	if v, ok := stats["draft_count"].(int64); ok {
		out.Body.DraftCount = v
	}
	if v, ok := stats["total_likes"].(int64); ok {
		out.Body.TotalLikes = v
	}
	if v, ok := stats["total_views"].(int64); ok {
		out.Body.TotalViews = v
	}
	if v, ok := stats["total_bookmarks"].(int64); ok {
		out.Body.TotalBookmarks = v
	}
	if v, ok := stats["emotion_distribution"].(map[string]int); ok {
		out.Body.EmotionDist = v
	}
	return out, nil
}
