package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-write-api/internal/features/admin/usecase"
)

// AdminHandler handles admin routes.
type AdminHandler struct{ uc usecase.AdminUseCase }

func NewAdminHandler(uc usecase.AdminUseCase) *AdminHandler { return &AdminHandler{uc: uc} }

// ── Input / Output types ────────────────────────────────────────────────────

type PaginatedInput struct {
	Page  int    `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	Q     string `query:"q" required:"false" doc:"Search query"`
}

type GetByIDInput struct {
	ID string `path:"id" format:"uuid" doc:"Resource UUID"`
}

type AdminStatsOutput struct {
	Body domain.DashboardStats
}

// ── Users ────────────────────────────────────────────────────────────────────

type ListUsersOutput struct {
	Body domain.PaginatedResponse[domain.AdminUser]
}

type GetUserOutput struct {
	Body domain.AdminUser
}

type CreateUserInput struct {
	Body struct {
		Username string `json:"username" minLength:"3" maxLength:"30" doc:"Username"`
		Email    string `json:"email" format:"email" doc:"Email address"`
		Password string `json:"password" minLength:"6" doc:"Password"`
		Role     string `json:"role" enum:"user,admin" doc:"User role"`
	}
}

type UpdateUserInput struct {
	ID   string `path:"id" format:"uuid" doc:"User UUID"`
	Body struct {
		Username  string `json:"username,omitempty" minLength:"3" maxLength:"30" required:"false" doc:"Username"`
		Email     string `json:"email,omitempty" required:"false" doc:"Email"`
		Bio       string `json:"bio,omitempty" required:"false" doc:"Bio"`
		AvatarURL string `json:"avatar_url,omitempty" required:"false" doc:"Avatar URL"`
		IsActive  *bool  `json:"is_active,omitempty" required:"false" doc:"Active status"`
	}
}

type ChangeRoleInput struct {
	ID   string `path:"id" format:"uuid" doc:"User UUID"`
	Body struct {
		Role string `json:"role" enum:"user,admin" doc:"New role"`
	}
}

// ── Poems ────────────────────────────────────────────────────────────────────

type ListPoemsInput struct {
	Page   int    `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	Q      string `query:"q" required:"false" doc:"Search query"`
	Status string `query:"status" required:"false" enum:"published,draft,removed" doc:"Filter by status"`
}

type ListPoemsOutput struct {
	Body domain.PaginatedResponse[domain.AdminPoem]
}

type GetPoemOutput struct {
	Body domain.AdminPoem
}

type AdminUpdatePoemInput struct {
	ID   string `path:"id" format:"uuid" doc:"Poem UUID"`
	Body struct {
		Title   string `json:"title,omitempty" minLength:"1" maxLength:"200" required:"false" doc:"Updated title"`
		Content string `json:"content,omitempty" minLength:"1" required:"false" doc:"Updated content"`
		Status  string `json:"status,omitempty" enum:"published,draft,removed" required:"false" doc:"Updated status"`
	}
}

type ChangeStatusInput struct {
	ID   string `path:"id" format:"uuid" doc:"Poem UUID"`
	Body struct {
		Status string `json:"status" enum:"published,draft,removed" doc:"New status"`
	}
}

// ── Associations ─────────────────────────────────────────────────────────────

type AssociationFilterInput struct {
	Page   int    `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	PoemID string `query:"poem_id" required:"false" doc:"Filter by poem UUID"`
	UserID string `query:"user_id" required:"false" doc:"Filter by user UUID"`
}

type ListLikesOutput struct {
	Body domain.PaginatedResponse[domain.AdminLike]
}
type ListBookmarksOutput struct {
	Body domain.PaginatedResponse[domain.AdminBookmark]
}
type ListEmotionsOutput struct {
	Body domain.PaginatedResponse[domain.AdminEmotion]
}

// ── Emotion Catalog ──────────────────────────────────────────────────────────

type EmotionCatalogListOutput struct {
	Body []*domain.EmotionCatalog
}

type CreateEmotionCatalogInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"50" doc:"Emotion name"`
		Emoji       string `json:"emoji" required:"false" doc:"Emoji representation"`
		Description string `json:"description" required:"false" doc:"Description"`
	}
}

type UpdateEmotionCatalogInput struct {
	ID   string `path:"id" format:"uuid" doc:"Emotion catalog UUID"`
	Body struct {
		Name        string `json:"name,omitempty" minLength:"1" maxLength:"50" required:"false" doc:"Updated name"`
		Emoji       string `json:"emoji,omitempty" required:"false" doc:"Updated emoji"`
		Description string `json:"description,omitempty" required:"false" doc:"Updated description"`
	}
}

type AdminMessageOutput struct {
	Body struct {
		Message string `json:"message" doc:"Result message"`
	}
}

// ── Handlers ────────────────────────────────────────────────────────────────

func (h *AdminHandler) GetStats(ctx context.Context, _ *struct{}) (*AdminStatsOutput, error) {
	stats, err := h.uc.GetDashboardStats()
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &AdminStatsOutput{Body: *stats}, nil
}

// Users

func (h *AdminHandler) ListUsers(ctx context.Context, input *PaginatedInput) (*ListUsersOutput, error) {
	result, err := h.uc.ListUsers(input.Page, input.Limit, input.Q)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListUsersOutput{Body: *result}, nil
}

func (h *AdminHandler) GetUser(ctx context.Context, input *GetByIDInput) (*GetUserOutput, error) {
	user, err := h.uc.GetUser(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetUserOutput{Body: *user}, nil
}

func (h *AdminHandler) CreateUser(ctx context.Context, input *CreateUserInput) (*AdminMessageOutput, error) {
	req := &domain.CreateUserRequest{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
		Role:     input.Body.Role,
	}
	if err := h.uc.CreateUser(req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "user created"
	return out, nil
}

func (h *AdminHandler) UpdateUser(ctx context.Context, input *UpdateUserInput) (*AdminMessageOutput, error) {
	req := &domain.UpdateUserAdminRequest{
		Username:  input.Body.Username,
		Email:     input.Body.Email,
		Bio:       input.Body.Bio,
		AvatarURL: input.Body.AvatarURL,
		IsActive:  input.Body.IsActive,
	}
	if err := h.uc.UpdateUser(input.ID, req); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "user updated"
	return out, nil
}

func (h *AdminHandler) ChangeUserRole(ctx context.Context, input *ChangeRoleInput) (*AdminMessageOutput, error) {
	if err := h.uc.ChangeUserRole(input.ID, input.Body.Role); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "role updated"
	return out, nil
}

func (h *AdminHandler) HardDeleteUser(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteUser(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

// Poems

func (h *AdminHandler) ListPoems(ctx context.Context, input *ListPoemsInput) (*ListPoemsOutput, error) {
	result, err := h.uc.ListPoems(input.Page, input.Limit, input.Q, input.Status)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListPoemsOutput{Body: *result}, nil
}

func (h *AdminHandler) GetPoem(ctx context.Context, input *GetByIDInput) (*GetPoemOutput, error) {
	poem, err := h.uc.GetPoem(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetPoemOutput{Body: *poem}, nil
}

func (h *AdminHandler) UpdatePoem(ctx context.Context, input *AdminUpdatePoemInput) (*AdminMessageOutput, error) {
	req := &domain.UpdatePoemAdminRequest{
		Title:   input.Body.Title,
		Content: input.Body.Content,
		Status:  input.Body.Status,
	}
	if err := h.uc.UpdatePoem(input.ID, req); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "poem updated"
	return out, nil
}

func (h *AdminHandler) ChangePoemStatus(ctx context.Context, input *ChangeStatusInput) (*AdminMessageOutput, error) {
	if err := h.uc.ChangePoemStatus(input.ID, input.Body.Status); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "status updated"
	return out, nil
}

func (h *AdminHandler) HardDeletePoem(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeletePoem(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

// Associations

func (h *AdminHandler) ListLikes(ctx context.Context, input *AssociationFilterInput) (*ListLikesOutput, error) {
	result, err := h.uc.ListLikes(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListLikesOutput{Body: *result}, nil
}

func (h *AdminHandler) HardDeleteLike(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteLike(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ListBookmarks(ctx context.Context, input *AssociationFilterInput) (*ListBookmarksOutput, error) {
	result, err := h.uc.ListBookmarks(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListBookmarksOutput{Body: *result}, nil
}

func (h *AdminHandler) HardDeleteBookmark(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteBookmark(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ListEmotions(ctx context.Context, input *AssociationFilterInput) (*ListEmotionsOutput, error) {
	result, err := h.uc.ListEmotions(input.Page, input.Limit, input.PoemID, input.UserID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ListEmotionsOutput{Body: *result}, nil
}

func (h *AdminHandler) HardDeleteEmotion(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.HardDeleteEmotion(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

// Emotion Catalog

func (h *AdminHandler) ListEmotionCatalog(ctx context.Context, _ *struct{}) (*EmotionCatalogListOutput, error) {
	items, err := h.uc.ListEmotionCatalog()
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &EmotionCatalogListOutput{Body: items}, nil
}

func (h *AdminHandler) CreateEmotionCatalog(ctx context.Context, input *CreateEmotionCatalogInput) (*AdminMessageOutput, error) {
	req := &domain.CreateEmotionCatalogRequest{
		Name:        input.Body.Name,
		Emoji:       input.Body.Emoji,
		Description: input.Body.Description,
	}
	if err := h.uc.CreateEmotionCatalog(req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion created"
	return out, nil
}

func (h *AdminHandler) UpdateEmotionCatalog(ctx context.Context, input *UpdateEmotionCatalogInput) (*AdminMessageOutput, error) {
	req := &domain.UpdateEmotionCatalogRequest{
		Name:        input.Body.Name,
		Emoji:       input.Body.Emoji,
		Description: input.Body.Description,
	}
	if err := h.uc.UpdateEmotionCatalog(input.ID, req); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "emotion updated"
	return out, nil
}

func (h *AdminHandler) DeleteEmotionCatalog(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.DeleteEmotionCatalog(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

// RegisterRoutes registers all admin Huma operations.
func RegisterAdminRoutes(api huma.API, h *AdminHandler, internalMW, requireAdminMW func(huma.Context, func(huma.Context))) {
	adminMW := huma.Middlewares{internalMW, requireAdminMW}

	huma.Register(api, huma.Operation{
		OperationID: "admin-stats",
		Summary:     "Get admin dashboard statistics",
		Method:      http.MethodGet,
		Path:        "/api/admin/stats",
		Tags:        []string{"Admin"},
		Middlewares:  adminMW,
	}, h.GetStats)

	// Users
	huma.Register(api, huma.Operation{
		OperationID: "admin-list-users",
		Summary:     "List all users (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/users",
		Tags:        []string{"Admin / Users"},
		Middlewares:  adminMW,
	}, h.ListUsers)

	huma.Register(api, huma.Operation{
		OperationID: "admin-get-user",
		Summary:     "Get a single user by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/users/{id}",
		Tags:        []string{"Admin / Users"},
		Middlewares:  adminMW,
	}, h.GetUser)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-user",
		Summary:       "Create a new user",
		Method:        http.MethodPost,
		Path:          "/api/admin/users",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Users"},
		Middlewares:    adminMW,
	}, h.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "admin-update-user",
		Summary:     "Update a user",
		Method:      http.MethodPut,
		Path:        "/api/admin/users/{id}",
		Tags:        []string{"Admin / Users"},
		Middlewares:  adminMW,
	}, h.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID: "admin-change-role",
		Summary:     "Change a user's role",
		Method:      http.MethodPatch,
		Path:        "/api/admin/users/{id}/role",
		Tags:        []string{"Admin / Users"},
		Middlewares:  adminMW,
	}, h.ChangeUserRole)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-user",
		Summary:       "Hard-delete a user",
		Method:        http.MethodDelete,
		Path:          "/api/admin/users/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Users"},
		Middlewares:    adminMW,
	}, h.HardDeleteUser)

	// Poems
	huma.Register(api, huma.Operation{
		OperationID: "admin-list-poems",
		Summary:     "List all poems (paginated, filterable)",
		Method:      http.MethodGet,
		Path:        "/api/admin/poems",
		Tags:        []string{"Admin / Poems"},
		Middlewares:  adminMW,
	}, h.ListPoems)

	huma.Register(api, huma.Operation{
		OperationID: "admin-get-poem",
		Summary:     "Get a single poem by ID",
		Method:      http.MethodGet,
		Path:        "/api/admin/poems/{id}",
		Tags:        []string{"Admin / Poems"},
		Middlewares:  adminMW,
	}, h.GetPoem)

	huma.Register(api, huma.Operation{
		OperationID: "admin-update-poem",
		Summary:     "Update a poem",
		Method:      http.MethodPut,
		Path:        "/api/admin/poems/{id}",
		Tags:        []string{"Admin / Poems"},
		Middlewares:  adminMW,
	}, h.UpdatePoem)

	huma.Register(api, huma.Operation{
		OperationID: "admin-change-poem-status",
		Summary:     "Change a poem's status",
		Method:      http.MethodPatch,
		Path:        "/api/admin/poems/{id}/status",
		Tags:        []string{"Admin / Poems"},
		Middlewares:  adminMW,
	}, h.ChangePoemStatus)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-poem",
		Summary:       "Hard-delete a poem",
		Method:        http.MethodDelete,
		Path:          "/api/admin/poems/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Poems"},
		Middlewares:    adminMW,
	}, h.HardDeletePoem)

	// Likes
	huma.Register(api, huma.Operation{
		OperationID: "admin-list-likes",
		Summary:     "List all likes (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/likes",
		Tags:        []string{"Admin / Likes"},
		Middlewares:  adminMW,
	}, h.ListLikes)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-like",
		Summary:       "Hard-delete a like",
		Method:        http.MethodDelete,
		Path:          "/api/admin/likes/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Likes"},
		Middlewares:    adminMW,
	}, h.HardDeleteLike)

	// Bookmarks
	huma.Register(api, huma.Operation{
		OperationID: "admin-list-bookmarks",
		Summary:     "List all bookmarks (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/bookmarks",
		Tags:        []string{"Admin / Bookmarks"},
		Middlewares:  adminMW,
	}, h.ListBookmarks)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-bookmark",
		Summary:       "Hard-delete a bookmark",
		Method:        http.MethodDelete,
		Path:          "/api/admin/bookmarks/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Bookmarks"},
		Middlewares:    adminMW,
	}, h.HardDeleteBookmark)

	// Emotions
	huma.Register(api, huma.Operation{
		OperationID: "admin-list-emotions",
		Summary:     "List all emotion tags (paginated)",
		Method:      http.MethodGet,
		Path:        "/api/admin/emotions",
		Tags:        []string{"Admin / Emotions"},
		Middlewares:  adminMW,
	}, h.ListEmotions)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-emotion",
		Summary:       "Hard-delete an emotion tag",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotions/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotions"},
		Middlewares:    adminMW,
	}, h.HardDeleteEmotion)

	// Emotion Catalog
	huma.Register(api, huma.Operation{
		OperationID: "admin-list-emotion-catalog",
		Summary:     "List all predefined emotions",
		Method:      http.MethodGet,
		Path:        "/api/admin/emotion-catalog",
		Tags:        []string{"Admin / Emotion Catalog"},
		Middlewares:  adminMW,
	}, h.ListEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-create-emotion-catalog",
		Summary:       "Create a new emotion in the catalog",
		Method:        http.MethodPost,
		Path:          "/api/admin/emotion-catalog",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin / Emotion Catalog"},
		Middlewares:    adminMW,
	}, h.CreateEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID: "admin-update-emotion-catalog",
		Summary:     "Update a catalog emotion",
		Method:      http.MethodPut,
		Path:        "/api/admin/emotion-catalog/{id}",
		Tags:        []string{"Admin / Emotion Catalog"},
		Middlewares:  adminMW,
	}, h.UpdateEmotionCatalog)

	huma.Register(api, huma.Operation{
		OperationID:   "admin-delete-emotion-catalog",
		Summary:       "Delete a catalog emotion",
		Method:        http.MethodDelete,
		Path:          "/api/admin/emotion-catalog/{id}",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin / Emotion Catalog"},
		Middlewares:    adminMW,
	}, h.DeleteEmotionCatalog)
}
