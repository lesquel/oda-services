package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error) {
	var tags []domain.EmotionTag
	db := r.db.Unscoped()
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.EmotionTag{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&tags).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminEmotion, len(tags))
	for i, t := range tags {
		items[i] = toAdminEmotion(t)
	}
	return &domain.PaginatedResponse[domain.AdminEmotion]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetEmotion(id string) (*domain.AdminEmotion, error) {
	var t domain.EmotionTag
	if err := r.db.Unscoped().First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	ae := toAdminEmotion(t)
	return &ae, nil
}

func (r *adminRepo) CreateEmotion(req *domain.CreateEmotionTagRequest) error {
	t := &domain.EmotionTag{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		PoemID:    req.PoemID,
		EmotionID: req.EmotionID,
	}
	return r.db.Create(t).Error
}

func (r *adminRepo) SoftDeleteEmotion(id string) error {
	return r.db.Delete(&domain.EmotionTag{}, "id = ?", id).Error
}

func (r *adminRepo) RestoreEmotion(id string) error {
	return r.db.Unscoped().Model(&domain.EmotionTag{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *adminRepo) PermanentDeleteEmotion(id string) error {
	return r.db.Unscoped().Delete(&domain.EmotionTag{}, "id = ?", id).Error
}
