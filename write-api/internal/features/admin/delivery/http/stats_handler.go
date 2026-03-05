package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lesquel/oda-shared/domain"
)

type AdminStatsOutput struct {
	Body domain.DashboardStats
}

func (h *AdminHandler) GetStats(ctx context.Context, _ *struct{}) (*AdminStatsOutput, error) {
	stats, err := h.uc.GetDashboardStats()
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return &AdminStatsOutput{Body: *stats}, nil
}
