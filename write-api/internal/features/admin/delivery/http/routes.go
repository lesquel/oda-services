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
		Summary:       "Soft-delete a user",
		Method:        http.MethodDelete,
		Path:          "/api/admin/users/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Users"},
		Middlewares:   adminMW,
	}, h.SoftDeleteUser)

	huma.Register(api, huma.Operation{
		OperationID: "admin-restore-user",
		Summary:     "Restore a soft-deleted user",
		Method:      http.MethodPost,
		Path:        "/api/admin/users/{id}/restore",
		Tags:        []string{"Admin / Users"},
		Middlewares: adminMW,
	}, h.RestoreUser)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-permanent-delete-user",
		Summary:       "Permanently delete a user",
		Method:        http.MethodDelete,
		Path:          "/api/admin/users/{id}/permanent",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Users"},
		Middlewares:   adminMW,
	}, h.PermanentDeleteUser)

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
		Summary:       "Soft-delete a poem",
		Method:        http.MethodDelete,
		Path:          "/api/admin/poems/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Poems"},
		Middlewares:   adminMW,
	}, h.SoftDeletePoem)

	huma.Register(api, huma.Operation{
		OperationID: "admin-restore-poem",
		Summary:     "Restore a soft-deleted poem",
		Method:      http.MethodPost,
		Path:        "/api/admin/poems/{id}/restore",
		Tags:        []string{"Admin / Poems"},
		Middlewares: adminMW,
	}, h.RestorePoem)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-permanent-delete-poem",
		Summary:       "Permanently delete a poem",
		Method:        http.MethodDelete,
		Path:          "/api/admin/poems/{id}/permanent",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Poems"},
		Middlewares:   adminMW,
	}, h.PermanentDeletePoem)

	// ── Moderation ───────────────────────────────────────────────────────

	huma.Register(api, huma.Operation{
		OperationID: "admin-moderation-queue",
		Summary:     "List poems pending moderation",
		Method:      http.MethodGet,
		Path:        "/api/admin/moderation/queue",
		Tags:        []string{"Admin / Moderation"},
		Middlewares: adminMW,
	}, h.ListModerationQueue)

	huma.Register(api, huma.Operation{
		OperationID: "admin-moderation-logs",
		Summary:     "Get moderation logs for a poem",
		Method:      http.MethodGet,
		Path:        "/api/admin/moderation/poems/{id}/logs",
		Tags:        []string{"Admin / Moderation"},
		Middlewares: adminMW,
	}, h.GetModerationLogs)

	huma.Register(api, huma.Operation{
		OperationID: "admin-moderation-action",
		Summary:     "Approve or reject a poem",
		Method:      http.MethodPost,
		Path:        "/api/admin/moderation/poems/{id}/action",
		Tags:        []string{"Admin / Moderation"},
		Middlewares: adminMW,
	}, h.ModerationAction)

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
		OperationID: "admin-get-like",
		Summary:     "Get a single like by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/likes/{id}",
		Tags:        []string{"Admin / Likes"},
		Middlewares: adminMW,
	}, h.GetLike)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-like",
		Summary:       "Create a like",
		Method:        http.MethodPost,
		Path:          "/api/admin/likes",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Likes"},
		Middlewares:   adminMW,
	}, h.CreateLike)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-like",
		Summary:       "Soft-delete a like",
		Method:        http.MethodDelete,
		Path:          "/api/admin/likes/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Likes"},
		Middlewares:   adminMW,
	}, h.SoftDeleteLike)

	huma.Register(api, huma.Operation{
		OperationID: "admin-restore-like",
		Summary:     "Restore a soft-deleted like",
		Method:      http.MethodPost,
		Path:        "/api/admin/likes/{id}/restore",
		Tags:        []string{"Admin / Likes"},
		Middlewares: adminMW,
	}, h.RestoreLike)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-permanent-delete-like",
		Summary:       "Permanently delete a like",
		Method:        http.MethodDelete,
		Path:          "/api/admin/likes/{id}/permanent",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Likes"},
		Middlewares:   adminMW,
	}, h.PermanentDeleteLike)

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
		OperationID: "admin-get-bookmark",
		Summary:     "Get a single bookmark by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/bookmarks/{id}",
		Tags:        []string{"Admin / Bookmarks"},
		Middlewares: adminMW,
	}, h.GetBookmark)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-bookmark",
		Summary:       "Create a bookmark",
		Method:        http.MethodPost,
		Path:          "/api/admin/bookmarks",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Bookmarks"},
		Middlewares:   adminMW,
	}, h.CreateBookmark)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-bookmark",
		Summary:       "Soft-delete a bookmark",
		Method:        http.MethodDelete,
		Path:          "/api/admin/bookmarks/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Bookmarks"},
		Middlewares:   adminMW,
	}, h.SoftDeleteBookmark)

	huma.Register(api, huma.Operation{
		OperationID: "admin-restore-bookmark",
		Summary:     "Restore a soft-deleted bookmark",
		Method:      http.MethodPost,
		Path:        "/api/admin/bookmarks/{id}/restore",
		Tags:        []string{"Admin / Bookmarks"},
		Middlewares: adminMW,
	}, h.RestoreBookmark)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-permanent-delete-bookmark",
		Summary:       "Permanently delete a bookmark",
		Method:        http.MethodDelete,
		Path:          "/api/admin/bookmarks/{id}/permanent",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Bookmarks"},
		Middlewares:   adminMW,
	}, h.PermanentDeleteBookmark)

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
		OperationID: "admin-get-emotion",
		Summary:     "Get a single emotion tag by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/emotions/{id}",
		Tags:        []string{"Admin / Emotions"},
		Middlewares: adminMW,
	}, h.GetEmotion)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-emotion",
		Summary:       "Create an emotion tag",
		Method:        http.MethodPost,
		Path:          "/api/admin/emotions",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Emotions"},
		Middlewares:   adminMW,
	}, h.CreateEmotion)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-emotion",
		Summary:       "Soft-delete an emotion tag",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotions/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotions"},
		Middlewares:   adminMW,
	}, h.SoftDeleteEmotion)

	huma.Register(api, huma.Operation{
		OperationID: "admin-restore-emotion",
		Summary:     "Restore a soft-deleted emotion tag",
		Method:      http.MethodPost,
		Path:        "/api/admin/emotions/{id}/restore",
		Tags:        []string{"Admin / Emotions"},
		Middlewares: adminMW,
	}, h.RestoreEmotion)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-permanent-delete-emotion",
		Summary:       "Permanently delete an emotion tag",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotions/{id}/permanent",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotions"},
		Middlewares:   adminMW,
	}, h.PermanentDeleteEmotion)

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
		OperationID: "admin-get-emotion-catalog",
		Summary:     "Get a single catalog emotion by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/emotion-catalog/{id}",
		Tags:        []string{"Admin / Emotion Catalog"},
		Middlewares: adminMW,
	}, h.GetEmotionCatalog)

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
		Summary:       "Soft-delete a catalog emotion",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotion-catalog/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotion Catalog"},
		Middlewares:   adminMW,
	}, h.SoftDeleteEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID: "admin-restore-emotion-catalog",
		Summary:     "Restore a soft-deleted catalog emotion",
		Method:      http.MethodPost,
		Path:        "/api/admin/emotion-catalog/{id}/restore",
		Tags:        []string{"Admin / Emotion Catalog"},
		Middlewares: adminMW,
	}, h.RestoreEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-permanent-delete-emotion-catalog",
		Summary:       "Permanently delete a catalog emotion",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotion-catalog/{id}/permanent",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotion Catalog"},
		Middlewares:   adminMW,
	}, h.PermanentDeleteEmotionCatalog)
}
