package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterAuthRoutes registers all auth-related Huma operations.
func RegisterAuthRoutes(api huma.API, h *AuthHandler, internalMW, requireUserMW func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Summary:       "Create a new user account",
		Method:        http.MethodPost,
		Path:          "/api/auth/register",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Auth"},
		Middlewares:   huma.Middlewares{internalMW},
	}, h.Register)

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Summary:     "Authenticate and get tokens",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{internalMW},
	}, h.Login)

	huma.Register(api, huma.Operation{
		OperationID: "refresh-token",
		Summary:     "Refresh an access token",
		Method:      http.MethodPost,
		Path:        "/api/auth/refresh",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{internalMW},
	}, h.Refresh)

	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Summary:       "Revoke a refresh token",
		Method:        http.MethodPost,
		Path:          "/api/auth/logout",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Auth"},
		Middlewares:   huma.Middlewares{internalMW, requireUserMW},
	}, h.Logout)

	huma.Register(api, huma.Operation{
		OperationID: "get-profile",
		Summary:     "Get the authenticated user's profile",
		Method:      http.MethodGet,
		Path:        "/api/me",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{internalMW, requireUserMW},
	}, h.GetProfile)

	huma.Register(api, huma.Operation{
		OperationID: "get-profile-legacy",
		Summary:     "Get the authenticated user's profile (legacy)",
		Method:      http.MethodGet,
		Path:        "/api/auth/profile",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{internalMW, requireUserMW},
	}, h.GetProfile)

	huma.Register(api, huma.Operation{
		OperationID: "update-profile",
		Summary:     "Update the authenticated user's profile",
		Method:      http.MethodPut,
		Path:        "/api/auth/profile",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{internalMW, requireUserMW},
	}, h.UpdateProfile)

	huma.Register(api, huma.Operation{
		OperationID: "change-password",
		Summary:     "Change the authenticated user's password",
		Method:      http.MethodPost,
		Path:        "/api/auth/change-password",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{internalMW, requireUserMW},
	}, h.ChangePassword)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-profile",
		Summary:     "Get a user's public profile by username",
		Method:      http.MethodGet,
		Path:        "/api/users/{username}",
		Tags:        []string{"Users"},
		Middlewares: huma.Middlewares{internalMW},
	}, h.GetPublicProfile)

	huma.Register(api, huma.Operation{
		OperationID: "search-users",
		Summary:     "Search users by username or email",
		Method:      http.MethodGet,
		Path:        "/api/users/search",
		Tags:        []string{"Users"},
		Middlewares: huma.Middlewares{internalMW},
	}, h.SearchUsers)
}
