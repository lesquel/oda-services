package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	var poems []domain.Poem
	db := r.db.Unscoped().Preload("Author")
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("title ILIKE ? OR content ILIKE ?", like, like)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	db.Model(&domain.Poem{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&poems).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminPoem, len(poems))
	for i, p := range poems {
		items[i] = toAdminPoem(p)
	}
	return &domain.PaginatedResponse[domain.AdminPoem]{
		Items: items, TotalCount: total,
		Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetPoem(id string) (*domain.AdminPoem, error) {
	var p domain.Poem
	if err := r.db.Unscoped().Preload("Author").First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	ap := toAdminPoem(p)
	return &ap, nil
}

func (r *adminRepo) UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error {
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) ChangePoemStatus(id, status string) error {
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).Update("status", status).Error
}

func (r *adminRepo) SoftDeletePoem(id string) error {
	return r.db.Delete(&domain.Poem{}, "id = ?", id).Error
}

func (r *adminRepo) RestorePoem(id string) error {
	return r.db.Unscoped().Model(&domain.Poem{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *adminRepo) PermanentDeletePoem(id string) error {
	return r.db.Unscoped().Delete(&domain.Poem{}, "id = ?", id).Error
}

// -- Moderation ---------------------------------------------------------------

func (r *adminRepo) ListModerationQueue(page, limit int) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	var poems []domain.Poem
	db := r.db.Preload("Author").Where("moderation_status = ? OR status = ?", "pending", "pending_review")
	var total int64
	db.Model(&domain.Poem{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Order("created_at ASC").Find(&poems).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminPoem, len(poems))
	for i, p := range poems {
		items[i] = toAdminPoem(p)
	}
	return &domain.PaginatedResponse[domain.AdminPoem]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetModerationLogs(poemID string) ([]domain.AdminModerationLog, error) {
	var logs []domain.ModerationLog
	if err := r.db.Where("poem_id = ?", poemID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminModerationLog, len(logs))
	for i, l := range logs {
		items[i] = domain.AdminModerationLog(l)
	}
	return items, nil
}

func (r *adminRepo) ModerationAction(poemID, action, reason, adminID string) error {
	now := time.Now()
	status := "approved"
	poemStatus := "published"
	if action == "reject" {
		status = "rejected"
		poemStatus = "rejected"
	}

	// Update poem
	updates := map[string]interface{}{
		"status":            poemStatus,
		"moderation_status": status,
		"moderation_reason": reason,
		"moderated_at":      now,
		"moderated_by":      "admin:" + adminID,
	}
	if err := r.db.Model(&domain.Poem{}).Where("id = ?", poemID).Updates(updates).Error; err != nil {
		return err
	}

	// Create log entry
	log := &domain.ModerationLog{
		ID:       uuid.NewString(),
		PoemID:   poemID,
		Status:   status,
		Reason:   reason,
		Provider: "admin",
		Model:    adminID,
	}
	return r.db.Create(log).Error
}
