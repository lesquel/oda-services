package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

// ── Poem types ──────────────────────────────────────────────────────────────

type ListPoemsInput struct {
	Page   int    `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	Q      string `query:"q" required:"false" doc:"Search query"`
	Status string `query:"status" required:"false" enum:"published,draft,removed,pending_review,rejected" doc:"Filter by status"`
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
		Status  string `json:"status,omitempty" enum:"published,draft,removed,pending_review,rejected" required:"false" doc:"Updated status"`
	}
}

type ChangeStatusInput struct {
	ID   string `path:"id" format:"uuid" doc:"Poem UUID"`
	Body struct {
		Status string `json:"status" enum:"published,draft,removed,pending_review,rejected" doc:"New status"`
	}
}

// ── Poem handlers ───────────────────────────────────────────────────────────

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

func (h *AdminHandler) SoftDeletePoem(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.SoftDeletePoem(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) RestorePoem(ctx context.Context, input *GetByIDInput) (*AdminMessageOutput, error) {
	if err := h.uc.RestorePoem(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "poem restored"
	return out, nil
}

func (h *AdminHandler) PermanentDeletePoem(ctx context.Context, input *GetByIDInput) (*struct{}, error) {
	if err := h.uc.PermanentDeletePoem(input.ID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

// ── Moderation types ────────────────────────────────────────────────────────

type ModerationQueueInput struct {
	Page  int `query:"page" default:"1" minimum:"1" doc:"Page number"`
	Limit int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
}

type ModerationQueueOutput struct {
	Body domain.PaginatedResponse[domain.AdminPoem]
}

type ModerationLogsOutput struct {
	Body struct {
		Logs []domain.AdminModerationLog `json:"logs"`
	}
}

type ModerationActionInput struct {
	ID   string `path:"id" format:"uuid" doc:"Poem UUID"`
	Body struct {
		Action string `json:"action" enum:"approve,reject" doc:"Moderation action"`
		Reason string `json:"reason" required:"false" doc:"Reason for the action"`
	}
}

// ── Moderation handlers ─────────────────────────────────────────────────────

func (h *AdminHandler) ListModerationQueue(ctx context.Context, input *ModerationQueueInput) (*ModerationQueueOutput, error) {
	result, err := h.uc.ListModerationQueue(input.Page, input.Limit)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &ModerationQueueOutput{Body: *result}, nil
}

func (h *AdminHandler) GetModerationLogs(ctx context.Context, input *GetByIDInput) (*ModerationLogsOutput, error) {
	logs, err := h.uc.GetModerationLogs(input.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &ModerationLogsOutput{}
	out.Body.Logs = logs
	return out, nil
}

func (h *AdminHandler) ModerationAction(ctx context.Context, input *ModerationActionInput) (*AdminMessageOutput, error) {
	adminID := "" // TODO: extract from auth context
	if err := h.uc.ModerationAction(input.ID, input.Body.Action, input.Body.Reason, adminID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	out := &AdminMessageOutput{}
	out.Body.Message = "moderation action applied"
	return out, nil
}
