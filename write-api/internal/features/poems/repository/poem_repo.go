package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

type poemRepo struct{ db *gorm.DB }

func NewPoemRepository(db *gorm.DB) domain.PoemRepository { return &poemRepo{db: db} }

func (r *poemRepo) Create(poem *domain.Poem) error {
	if poem.ID == "" {
		poem.ID = uuid.NewString()
	}
	return r.db.Create(poem).Error
}

func (r *poemRepo) FindByID(id string) (*domain.Poem, error) {
	var poem domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").First(&poem, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &poem, nil
}

func (r *poemRepo) Update(poem *domain.Poem) error {
	return r.db.Omit("Author", "EmotionTags").Save(poem).Error
}

func (r *poemRepo) Delete(id string) error {
	return r.db.Delete(&domain.Poem{}, "id = ?", id).Error
}

func (r *poemRepo) GetFeed(limit, offset int) ([]*domain.Poem, error) {
	var poems []*domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, err
}

func (r *poemRepo) GetUserPoems(userID string, limit, offset int) ([]*domain.Poem, error) {
	var poems []*domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, err
}

func (r *poemRepo) Search(query string, limit, offset int) ([]*domain.Poem, error) {
	var poems []*domain.Poem
	like := "%" + query + "%"
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("status = ? AND (title ILIKE ? OR content ILIKE ?)", "published", like, like).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, err
}

func (r *poemRepo) IncrementViews(id string) error {
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
}

func (r *poemRepo) GetStats(poemID string) (map[string]interface{}, error) {
	var poem domain.Poem
	if err := r.db.Select("id, likes_count, views_count").First(&poem, "id = ?", poemID).Error; err != nil {
		return nil, err
	}
	var emotionCount int64
	r.db.Model(&domain.EmotionTag{}).Where("poem_id = ?", poemID).Count(&emotionCount)
	return map[string]interface{}{
		"likes_count":   poem.LikesCount,
		"views_count":   poem.ViewsCount,
		"emotion_count": emotionCount,
	}, nil
}
