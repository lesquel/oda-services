package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Find(&items).Error
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

func (r *adminRepo) DeleteEmotionCatalog(id string) error {
	return r.db.Delete(&domain.EmotionCatalog{}, "id = ?", id).Error
}
