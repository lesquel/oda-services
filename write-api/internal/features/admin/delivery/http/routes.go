package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterAdminRoutes registers all admin Huma operations.
func RegisterAdminRoutes(api huma.API, h *AdminHandler, internalMW, requireAdminMW func(huma.Context, func(huma.Context))) {
	adminMW := huma.Middlewares{internalMW, requireAdminMW}

	// ── Stats ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-stats",
		Summary:     "Get admin dashboard statistics",
		Method:      http.MethodGet,
		Path:        "/api/admin/stats",
		Tags:        []string{"Admin"},
		Middlewares: adminMW,
	}, h.GetStats)

	// ── Users ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-list-users",
		Summary:     "List all users (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/users",
		Tags:        []string{"Admin / Users"},
		Middlewares: adminMW,
	}, h.ListUsers)

	huma.Register(api, huma.Operation{
		OperationID: "admin-get-user",
		Summary:     "Get a single user by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/users/{id}",
		Tags:        []string{"Admin / Users"},
		Middlewares: adminMW,
	}, h.GetUser)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-user",
		Summary:       "Create a new user",
		Method:        http.MethodPost,
		Path:          "/api/admin/users",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Users"},
		Middlewares:   adminMW,
	}, h.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "admin-update-user",
		Summary:     "Update a user",
		Method:      http.MethodPut,
		Path:        "/api/admin/users/{id}",
		Tags:        []string{"Admin / Users"},
		Middlewares: adminMW,
	}, h.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID: "admin-change-role",
		Summary:     "Change a user's role",
		Method:      http.MethodPatch,
		Path:        "/api/admin/users/{id}/role",
		Tags:        []string{"Admin / Users"},
		Middlewares: adminMW,
	}, h.ChangeUserRole)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-user",
		Summary:       "Hard-delete a user",
		Method:        http.MethodDelete,
		Path:          "/api/admin/users/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Users"},
		Middlewares:   adminMW,
	}, h.HardDeleteUser)

	// ── Poems ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-list-poems",
		Summary:     "List all poems (paginated, filterable)",
		Method:      http.MethodGet,
		Path:        "/api/admin/poems",
		Tags:        []string{"Admin / Poems"},
		Middlewares: adminMW,
	}, h.ListPoems)

	huma.Register(api, huma.Operation{
		OperationID: "admin-get-poem",
		Summary:     "Get a single poem by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/poems/{id}",
		Tags:        []string{"Admin / Poems"},
		Middlewares: adminMW,
	}, h.GetPoem)

	huma.Register(api, huma.Operation{
		OperationID: "admin-update-poem",
		Summary:     "Update a poem",
		Method:      http.MethodPut,
		Path:        "/api/admin/poems/{id}",
		Tags:        []string{"Admin / Poems"},
		Middlewares: adminMW,
	}, h.UpdatePoem)

	huma.Register(api, huma.Operation{
		OperationID: "admin-change-poem-status",
		Summary:     "Change a poem's status",
		Method:      http.MethodPatch,
		Path:        "/api/admin/poems/{id}/status",
		Tags:        []string{"Admin / Poems"},
		Middlewares: adminMW,
	}, h.ChangePoemStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-poem",
		Summary:       "Hard-delete a poem",
		Method:        http.MethodDelete,
		Path:          "/api/admin/poems/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Poems"},
		Middlewares:   adminMW,
	}, h.HardDeletePoem)

	// ── Likes ────────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-list-likes",
		Summary:     "List all likes (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/likes",
		Tags:        []string{"Admin / Likes"},
		Middlewares: adminMW,
	}, h.ListLikes)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-like",
		Summary:       "Hard-delete a like",
		Method:        http.MethodDelete,
		Path:          "/api/admin/likes/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Likes"},
		Middlewares:   adminMW,
	}, h.HardDeleteLike)

	// ── Bookmarks ────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-list-bookmarks",
		Summary:     "List all bookmarks (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/bookmarks",
		Tags:        []string{"Admin / Bookmarks"},
		Middlewares: adminMW,
	}, h.ListBookmarks)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-bookmark",
		Summary:       "Hard-delete a bookmark",
		Method:        http.MethodDelete,
		Path:          "/api/admin/bookmarks/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Bookmarks"},
		Middlewares:   adminMW,
	}, h.HardDeleteBookmark)

	// ── Emotions ─────────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-list-emotions",
		Summary:     "List all emotion tags (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/emotions",
		Tags:        []string{"Admin / Emotions"},
		Middlewares: adminMW,
	}, h.ListEmotions)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-emotion",
		Summary:       "Hard-delete an emotion tag",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotions/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotions"},
		Middlewares:   adminMW,
	}, h.HardDeleteEmotion)

	// ── Emotion Catalog ──────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-list-emotion-catalog",
		Summary:     "List all predefined emotions",
		Method:      http.MethodGet,
		Path:        "/api/admin/emotion-catalog",
		Tags:        []string{"Admin / Emotion Catalog"},
		Middlewares: adminMW,
	}, h.ListEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-emotion-catalog",
		Summary:       "Create a new emotion in the catalog",
		Method:        http.MethodPost,
		Path:          "/api/admin/emotion-catalog",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Emotion Catalog"},
		Middlewares:   adminMW,
	}, h.CreateEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID: "admin-update-emotion-catalog",
		Summary:     "Update a catalog emotion",
		Method:      http.MethodPut,
		Path:        "/api/admin/emotion-catalog/{id}",
		Tags:        []string{"Admin / Emotion Catalog"},
		Middlewares: adminMW,
	}, h.UpdateEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-emotion-catalog",
		Summary:       "Delete a catalog emotion",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotion-catalog/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotion Catalog"},
		Middlewares:   adminMW,
	}, h.DeleteEmotionCatalog)
}
