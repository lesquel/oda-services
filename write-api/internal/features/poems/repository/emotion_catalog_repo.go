package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

type emotionCatalogRepo struct{ db *gorm.DB }

func NewEmotionCatalogRepository(db *gorm.DB) domain.EmotionCatalogRepository {
	return &emotionCatalogRepo{db: db}
}

func (r *emotionCatalogRepo) FindAll() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Find(&items).Error
}

func (r *emotionCatalogRepo) FindByID(id string) (*domain.EmotionCatalog, error) {
	var item domain.EmotionCatalog
	return &item, r.db.First(&item, "id = ?", id).Error
}

func (r *emotionCatalogRepo) FindByName(name string) (*domain.EmotionCatalog, error) {
	var item domain.EmotionCatalog
	return &item, r.db.First(&item, "name = ?", name).Error
}

func (r *emotionCatalogRepo) Create(e *domain.EmotionCatalog) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return r.db.Create(e).Error
}

func (r *emotionCatalogRepo) Update(e *domain.EmotionCatalog) error {
	return r.db.Save(e).Error
}

func (r *emotionCatalogRepo) Delete(id string) error {
	return r.db.Delete(&domain.EmotionCatalog{}, "id = ?", id).Error
}
