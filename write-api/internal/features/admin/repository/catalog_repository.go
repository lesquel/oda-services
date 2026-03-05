package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Unscoped().Find(&items).Error
}

func (r *adminRepo) GetEmotionCatalog(id string) (*domain.EmotionCatalog, error) {
	var item domain.EmotionCatalog
	if err := r.db.Unscoped().First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *adminRepo) CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error {
	item := &domain.EmotionCatalog{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Emoji:       req.Emoji,
		Description: req.Description,
	}
	return r.db.Create(item).Error
}

func (r *adminRepo) UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error {
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Emoji != "" {
		updates["emoji"] = req.Emoji
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	return r.db.Model(&domain.EmotionCatalog{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) SoftDeleteEmotionCatalog(id string) error {
	return r.db.Delete(&domain.EmotionCatalog{}, "id = ?", id).Error
}

func (r *adminRepo) RestoreEmotionCatalog(id string) error {
	return r.db.Unscoped().Model(&domain.EmotionCatalog{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *adminRepo) PermanentDeleteEmotionCatalog(id string) error {
	return r.db.Unscoped().Delete(&domain.EmotionCatalog{}, "id = ?", id).Error
}
