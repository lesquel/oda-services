package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── User types ──────────────────────────────────────────────────────────────

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

// ── User handlers ───────────────────────────────────────────────────────────

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

func (h *AdminHandler) SoftDeleteUser(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.SoftDeleteUser(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) RestoreUser(ctx context.Context, input *GetByIDInput) (*AdminMessageOutput, error) {
	if err := h.uc.RestoreUser(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "user restored"
	return out, nil
}

func (h *AdminHandler) PermanentDeleteUser(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.PermanentDeleteUser(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}
