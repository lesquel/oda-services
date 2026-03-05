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
	Body domain.User
}

type SearchUsersInput struct {
	Q string `query:"q" required:"false" doc:"Search query"`
}
type SearchUsersOutput struct {
	Body []*domain.User
}

// ── User handlers ───────────────────────────────────────────────────────────

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
