package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Auth types ──────────────────────────────────────────────────────────────

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

// ── Auth handlers ───────────────────────────────────────────────────────────

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
