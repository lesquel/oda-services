package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-write-api/internal/features/auth/usecase"
	"github.com/lesquel/oda-write-api/internal/middleware"
)

// AuthHandler handles auth and user mutation routes.
type AuthHandler struct{ uc usecase.AuthUseCase }

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler { return &AuthHandler{uc: uc} }

// ── Input / Output types ────────────────────────────────────────────────────

type RegisterInput struct {
	Body struct {
		Username string `json:"username" minLength:"3" maxLength:"30" doc:"Username"`
		Email    string `json:"email" format:"email" doc:"Email address"`
		Password string `json:"password" minLength:"6" doc:"Password"`
	}
}
type RegisterOutput struct {
	Body domain.AuthResponse
}

type LoginInput struct {
	Body struct {
		Email    string `json:"email" format:"email" doc:"Email address"`
		Password string `json:"password" doc:"Password"`
	}
}
type LoginOutput struct {
	Body domain.AuthResponse
}

type RefreshInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" doc:"Refresh token"`
	}
}
type RefreshOutput struct {
	Body domain.AuthResponse
}

type LogoutInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" required:"false" doc:"Refresh token to revoke"`
	}
}

type GetProfileOutput struct {
	Body domain.User
}

type UpdateProfileInput struct {
	Body struct {
		Username  string `json:"username,omitempty" minLength:"3" maxLength:"30" required:"false" doc:"New username"`
		Bio       string `json:"bio,omitempty" required:"false" doc:"Bio text"`
		AvatarURL string `json:"avatar_url,omitempty" required:"false" doc:"Avatar URL"`
	}
}
type UpdateProfileOutput struct {
	Body domain.User
}

type ChangePasswordInput struct {
	Body struct {
		OldPassword string `json:"old_password" doc:"Current password"`
		NewPassword string `json:"new_password" minLength:"6" doc:"New password"`
	}
}

type MessageOutput struct {
	Body struct {
		Message string `json:"message" doc:"Response message"`
	}
}

type GetPublicProfileInput struct {
	Username string `path:"username" doc:"Username to look up"`
}
type GetPublicProfileOutput struct {
	Body domain.User
}

type SearchUsersInput struct {
	Q string `query:"q" required:"false" doc:"Search query"`
}
type SearchUsersOutput struct {
	Body []*domain.User
}

// ── Handlers ────────────────────────────────────────────────────────────────

func (h *AuthHandler) Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	req := &domain.RegisterRequest{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
	}
	resp, err := h.uc.Register(req)
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	return &RegisterOutput{Body: *resp}, nil
}

func (h *AuthHandler) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	req := &domain.LoginRequest{
		Email:    input.Body.Email,
		Password: input.Body.Password,
	}
	resp, err := h.uc.Login(req)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, err.Error())
	}
	return &LoginOutput{Body: *resp}, nil
}

func (h *AuthHandler) Refresh(ctx context.Context, input *RefreshInput) (*RefreshOutput, error) {
	resp, err := h.uc.Refresh(input.Body.RefreshToken)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, err.Error())
	}
	return &RefreshOutput{Body: *resp}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, input *LogoutInput) (*struct{}, error) {
	_ = h.uc.Logout(input.Body.RefreshToken)
	return nil, nil
}

func (h *AuthHandler) GetProfile(ctx context.Context, input *struct{}) (*GetProfileOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	user, err := h.uc.GetProfile(userID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetProfileOutput{Body: *user}, nil
}

func (h *AuthHandler) UpdateProfile(ctx context.Context, input *UpdateProfileInput) (*UpdateProfileOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	req := &domain.UpdateProfileRequest{
		Username:  input.Body.Username,
		Bio:       input.Body.Bio,
		AvatarURL: input.Body.AvatarURL,
	}
	user, err := h.uc.UpdateProfile(userID, req)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &UpdateProfileOutput{Body: *user}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, input *ChangePasswordInput) (*MessageOutput, error) {
	userID, _ := middleware.GetUserID(ctx)
	req := &domain.ChangePasswordRequest{
		OldPassword: input.Body.OldPassword,
		NewPassword: input.Body.NewPassword,
	}
	if err := h.uc.ChangePassword(userID, req); err != nil {
		return nil, huma.NewError(http.StatusBadRequest, err.Error())
	}
	out := &MessageOutput{}
	out.Body.Message = "contraseña actualizada"
	return out, nil
}

func (h *AuthHandler) GetPublicProfile(ctx context.Context, input *GetPublicProfileInput) (*GetPublicProfileOutput, error) {
	profile, err := h.uc.GetPublicProfile(input.Username)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, err.Error())
	}
	return &GetPublicProfileOutput{Body: *profile}, nil
}

func (h *AuthHandler) SearchUsers(ctx context.Context, input *SearchUsersInput) (*SearchUsersOutput, error) {
	if input.Q == "" {
		return &SearchUsersOutput{Body: []*domain.User{}}, nil
	}
	profiles, err := h.uc.SearchUsers(input.Q)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &SearchUsersOutput{Body: profiles}, nil
}

// RegisterRoutes registers all auth-related Huma operations.
func RegisterAuthRoutes(api huma.API, h *AuthHandler, internalMW, requireUserMW func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Summary:       "Create a new user account",
		Method:        http.MethodPost,
		Path:          "/api/auth/register",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Auth"},
		Middlewares:    huma.Middlewares{internalMW},
	}, h.Register)

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Summary:     "Authenticate and get tokens",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Tags:        []string{"Auth"},
		Middlewares:  huma.Middlewares{internalMW},
	}, h.Login)

	huma.Register(api, huma.Operation{
		OperationID: "refresh-token",
		Summary:     "Refresh an access token",
		Method:      http.MethodPost,
		Path:        "/api/auth/refresh",
		Tags:        []string{"Auth"},
		Middlewares:  huma.Middlewares{internalMW},
	}, h.Refresh)

	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Summary:       "Revoke a refresh token",
		Method:        http.MethodPost,
		Path:          "/api/auth/logout",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Auth"},
		Middlewares:    huma.Middlewares{internalMW, requireUserMW},
	}, h.Logout)

	huma.Register(api, huma.Operation{
		OperationID: "get-profile",
		Summary:     "Get the authenticated user's profile",
		Method:      http.MethodGet,
		Path:        "/api/me",
		Tags:        []string{"Auth"},
		Middlewares:  huma.Middlewares{internalMW, requireUserMW},
	}, h.GetProfile)

	huma.Register(api, huma.Operation{
		OperationID: "get-profile-legacy",
		Summary:     "Get the authenticated user's profile (legacy)",
		Method:      http.MethodGet,
		Path:        "/api/auth/profile",
		Tags:        []string{"Auth"},
		Middlewares:  huma.Middlewares{internalMW, requireUserMW},
	}, h.GetProfile)

	huma.Register(api, huma.Operation{
		OperationID: "update-profile",
		Summary:     "Update the authenticated user's profile",
		Method:      http.MethodPut,
		Path:        "/api/auth/profile",
		Tags:        []string{"Auth"},
		Middlewares:  huma.Middlewares{internalMW, requireUserMW},
	}, h.UpdateProfile)

	huma.Register(api, huma.Operation{
		OperationID: "change-password",
		Summary:     "Change the authenticated user's password",
		Method:      http.MethodPost,
		Path:        "/api/auth/change-password",
		Tags:        []string{"Auth"},
		Middlewares:  huma.Middlewares{internalMW, requireUserMW},
	}, h.ChangePassword)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-profile",
		Summary:     "Get a user's public profile by username",
		Method:      http.MethodGet,
		Path:        "/api/users/{username}",
		Tags:        []string{"Users"},
		Middlewares:  huma.Middlewares{internalMW},
	}, h.GetPublicProfile)

	huma.Register(api, huma.Operation{
		OperationID: "search-users",
		Summary:     "Search users by username or email",
		Method:      http.MethodGet,
		Path:        "/api/users/search",
		Tags:        []string{"Users"},
		Middlewares:  huma.Middlewares{internalMW},
	}, h.SearchUsers)
}
