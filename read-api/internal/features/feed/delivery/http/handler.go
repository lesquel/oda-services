package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-read-api/internal/features/feed/usecase"
	"github.com/lesquel/oda-read-api/internal/middleware"
	"github.com/lesquel/oda-shared/domain"
)

// ReadHandler handles all read-side Huma operations.
type ReadHandler struct{ uc *usecase.ReadUseCase }

func NewReadHandler(uc *usecase.ReadUseCase) *ReadHandler {
	return &ReadHandler{uc: uc}
}

// ── Input / Output types ────────────────────────────────────────────────────

// -- Feed --

type GetFeedInput struct {
	Page  int `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit int `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type GetFeedOutput struct {
	Body []*domain.Poem
}

// -- Search poems --

type SearchPoemsInput struct {
	Q     string `query:"q" required:"false" doc:"Search query"`
	Page  int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type SearchPoemsOutput struct {
	Body []*domain.Poem
}

// -- Single poem --

type GetPoemInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem ID"`
}
type GetPoemOutput struct {
	Body *domain.Poem
}

// -- Poem stats --

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

// -- Emotion distribution --

type GetEmotionDistInput struct {
	ID string `path:"id" format:"uuid" doc:"Poem ID"`
}
type GetEmotionDistOutput struct {
	Body map[string]int
}

// -- User poems --

type GetUserPoemsInput struct {
	UserID string `path:"userID" format:"uuid" doc:"User ID"`
	Status string `query:"status" required:"false" enum:"published,draft" default:"published" doc:"Filter by poem status"`
	Page   int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type GetUserPoemsOutput struct {
	Body []*domain.Poem
}

// -- Public profile --

type GetPublicProfileInput struct {
	Username string `path:"username" doc:"Username to look up"`
}
type GetPublicProfileOutput struct {
	Body *domain.User
}

// -- Search users --

type SearchUsersInput struct {
	Q      string `query:"q" required:"false" doc:"Search query"`
	Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Max results"`
	Offset int    `query:"offset" minimum:"0" default:"0" doc:"Offset for pagination"`
}
type SearchUsersOutput struct {
	Body []*domain.User
}

// -- User stats --

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

// -- Bookmarks --

type GetBookmarksInput struct {
	Page  int `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit int `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}
type GetBookmarksOutput struct {
	Body []*domain.Poem
}

// -- Emotion catalog --

type GetEmotionCatalogOutput struct {
	Body []*domain.EmotionCatalog
}

// ── Handlers ────────────────────────────────────────────────────────────────

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

func (h *ReadHandler) GetUserBookmarks(ctx context.Context, input *GetBookmarksInput) (*GetBookmarksOutput, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, huma.NewError(http.StatusUnauthorized, "authentication required")
	}
	poems, _, err := h.uc.GetUserBookmarks(userID, input.Page, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	if poems == nil {
		poems = []*domain.Poem{}
	}
	return &GetBookmarksOutput{Body: poems}, nil
}

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

// ── Route registration ──────────────────────────────────────────────────────

func RegisterReadRoutes(api huma.API, h *ReadHandler, internalMW, requireUserMW func(huma.Context, func(huma.Context))) {
	authMW := huma.Middlewares{internalMW}
	userMW := huma.Middlewares{internalMW, requireUserMW}

	// ── Poems ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "get-feed",
		Summary:     "Get the public poem feed",
		Method:      http.MethodGet,
		Path:        "/api/poems/feed",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.GetFeed)

	huma.Register(api, huma.Operation{
		OperationID: "search-poems",
		Summary:     "Search poems by title or content",
		Method:      http.MethodGet,
		Path:        "/api/poems/search",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.SearchPoems)

	huma.Register(api, huma.Operation{
		OperationID: "get-poem",
		Summary:     "Get a single poem by ID",
		Method:      http.MethodGet,
		Path:        "/api/poems/{id}",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.GetPoem)

	huma.Register(api, huma.Operation{
		OperationID: "get-poem-stats",
		Summary:     "Get likes, views and emotion count for a poem",
		Method:      http.MethodGet,
		Path:        "/api/poems/{id}/stats",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.GetPoemStats)

	huma.Register(api, huma.Operation{
		OperationID: "get-emotion-distribution",
		Summary:     "Get emotion distribution for a poem",
		Method:      http.MethodGet,
		Path:        "/api/poems/{id}/emotions/distribution",
		Tags:        []string{"Poems"},
		Middlewares:  authMW,
	}, h.GetEmotionDistribution)

	// ── Users ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "search-users-read",
		Summary:     "Search users",
		Method:      http.MethodGet,
		Path:        "/api/users/search",
		Tags:        []string{"Users"},
		Middlewares:  authMW,
	}, h.SearchUsers)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-profile-read",
		Summary:     "Get a user's public profile",
		Method:      http.MethodGet,
		Path:        "/api/users/{username}",
		Tags:        []string{"Users"},
		Middlewares:  authMW,
	}, h.GetPublicProfile)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-poems",
		Summary:     "Get poems by a specific user",
		Method:      http.MethodGet,
		Path:        "/api/users/{userID}/poems",
		Tags:        []string{"Users"},
		Middlewares:  authMW,
	}, h.GetUserPoems)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-stats",
		Summary:     "Get aggregate stats for a user",
		Method:      http.MethodGet,
		Path:        "/api/users/{userID}/stats",
		Tags:        []string{"Users"},
		Middlewares:  authMW,
	}, h.GetUserStats)

	// ── Bookmarks ────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "get-bookmarks",
		Summary:     "Get the authenticated user's bookmarks",
		Method:      http.MethodGet,
		Path:        "/api/bookmarks",
		Tags:        []string{"Bookmarks"},
		Middlewares:  userMW,
	}, h.GetUserBookmarks)

	// ── Emotion catalog ──────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "get-emotion-catalog",
		Summary:     "List all available emotions",
		Method:      http.MethodGet,
		Path:        "/api/emotions",
		Tags:        []string{"Emotions"},
		Middlewares:  authMW,
	}, h.GetEmotionCatalog)
}
