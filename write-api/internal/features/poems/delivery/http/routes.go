package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterPoemRoutes registers all poem-related Huma operations.
func RegisterPoemRoutes(api huma.API, h *PoemHandler, internalMW, requireUserMW func(huma.Context, func(huma.Context))) {
	authMW := huma.Middlewares{internalMW, requireUserMW}

	huma.Register(api, huma.Operation{
		OperationID:   "create-poem",
		Summary:       "Create a new poem",
		Method:        http.MethodPost,
		Path:          "/api/poems",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Poems"},
		Middlewares:   authMW,
	}, h.CreatePoem)

	huma.Register(api, huma.Operation{
		OperationID: "update-poem",
		Summary:     "Update an existing poem",
		Method:      http.MethodPut,
		Path:        "/api/poems/{id}",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.UpdatePoem)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-poem",
		Summary:       "Delete a poem",
		Method:        http.MethodDelete,
		Path:          "/api/poems/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Poems"},
		Middlewares:   authMW,
	}, h.DeletePoem)

	huma.Register(api, huma.Operation{
		OperationID: "toggle-like",
		Summary:     "Like or unlike a poem",
		Method:      http.MethodPost,
		Path:        "/api/poems/{id}/like",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.ToggleLike)

	huma.Register(api, huma.Operation{
		OperationID: "toggle-bookmark",
		Summary:     "Bookmark or unbookmark a poem",
		Method:      http.MethodPost,
		Path:        "/api/poems/{id}/bookmark",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.ToggleBookmark)

	huma.Register(api, huma.Operation{
		OperationID: "tag-emotion",
		Summary:     "Tag an emotion on a poem",
		Method:      http.MethodPost,
		Path:        "/api/poems/{id}/emotions",
		Tags:        []string{"Poems"},
		Middlewares: authMW,
	}, h.TagEmotion)

	huma.Register(api, huma.Operation{
		OperationID:   "remove-emotion-tag",
		Summary:       "Remove your emotion tag from a poem",
		Method:        http.MethodDelete,
		Path:          "/api/poems/{id}/emotions/{emotionID}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Poems"},
		Middlewares:   authMW,
	}, h.RemoveEmotionTag)
}
