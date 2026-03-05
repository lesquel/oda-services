package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-write-api/internal/middleware"
)

// ── Profile types ───────────────────────────────────────────────────────────

type GetProfileOutput struct {
	Body domain.User
}

type UpdateProfileInput struct {
	Body struct {
		Username  string `json:"username,omitempty" minLength:"3" maxLength:"30" required:"false" doc:"New username"`
		Bio       string `json:"bio,omitempty" required:"false" doc:"Bio text"`
		AvatarURL string `json:"avatar_url,omitempty" required:"false" doc:"Avatar URL"`
		Website   string `json:"website,omitempty" required:"false" doc:"Website URL"`
		Instagram string `json:"instagram,omitempty" required:"false" doc:"Instagram handle"`
		Twitter   string `json:"twitter,omitempty" required:"false" doc:"X/Twitter handle"`
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

// ── Profile handlers ────────────────────────────────────────────────────────

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
		Website:   input.Body.Website,
		Instagram: input.Body.Instagram,
		Twitter:   input.Body.Twitter,
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
