package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterReadRoutes wires all read-side Huma operations.
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
		Middlewares: authMW,
	}, h.GetFeed)

	huma.Register(api, huma.Operation{
		OperationID: "search-poems",
		Summary:     "Search poems by title or content",
		Method:      http.MethodGet,
		Path:        "/api/poems/search",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.SearchPoems)

	huma.Register(api, huma.Operation{
		OperationID: "get-poem",
		Summary:     "Get a single poem by ID",
		Method:      http.MethodGet,
		Path:        "/api/poems/{id}",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.GetPoem)

	huma.Register(api, huma.Operation{
		OperationID: "get-poem-stats",
		Summary:     "Get likes, views and emotion count for a poem",
		Method:      http.MethodGet,
		Path:        "/api/poems/{id}/stats",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.GetPoemStats)

	huma.Register(api, huma.Operation{
		OperationID: "get-emotion-distribution",
		Summary:     "Get emotion distribution for a poem",
		Method:      http.MethodGet,
		Path:        "/api/poems/{id}/emotions/distribution",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.GetEmotionDistribution)

	// ── Users ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "search-users-read",
		Summary:     "Search users",
		Method:      http.MethodGet,
		Path:        "/api/users/search",
		Tags:        []string{"Users"},
		Middlewares: authMW,
	}, h.SearchUsers)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-profile-read",
		Summary:     "Get a user's public profile",
		Method:      http.MethodGet,
		Path:        "/api/users/{username}",
		Tags:        []string{"Users"},
		Middlewares: authMW,
	}, h.GetPublicProfile)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-poems",
		Summary:     "Get poems by a specific user",
		Method:      http.MethodGet,
		Path:        "/api/users/{userID}/poems",
		Tags:        []string{"Users"},
		Middlewares: authMW,
	}, h.GetUserPoems)

	huma.Register(api, huma.Operation{
		OperationID: "get-user-stats",
		Summary:     "Get aggregate stats for a user",
		Method:      http.MethodGet,
		Path:        "/api/users/{userID}/stats",
		Tags:        []string{"Users"},
		Middlewares: authMW,
	}, h.GetUserStats)

	// ── Bookmarks ────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "get-bookmarks",
		Summary:     "Get the authenticated user's bookmarks",
		Method:      http.MethodGet,
		Path:        "/api/bookmarks",
		Tags:        []string{"Bookmarks"},
		Middlewares: userMW,
	}, h.GetUserBookmarks)

	// ── Emotion catalog ──────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "get-emotion-catalog",
		Summary:     "List all available emotions",
		Method:      http.MethodGet,
		Path:        "/api/emotions",
		Tags:        []string{"Emotions"},
		Middlewares: authMW,
	}, h.GetEmotionCatalog)
}
