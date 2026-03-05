package repository

import "github.com/lesquel/oda-shared/domain"

// ── Feed ──────────────────────────────────────────────────────────────────────

func (r *ReadRepository) GetFeed(limit, offset int) ([]*domain.Poem, int64, error) {
	var poems []*domain.Poem
	var total int64
	r.db.Model(&domain.Poem{}).Where("status = ?", "published").Count(&total)
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, total, err
}

// ── Single poem ───────────────────────────────────────────────────────────────

func (r *ReadRepository) GetPoem(id string) (*domain.Poem, error) {
	var poem domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").
		First(&poem, "id = ? AND status = ?", id, "published").Error
	if err != nil {
		return nil, err
	}
	return &poem, nil
}

func (r *ReadRepository) SearchPoems(query string, limit, offset int) ([]*domain.Poem, int64, error) {
	var poems []*domain.Poem
	var total int64
	like := "%" + query + "%"
	db := r.db.Where("status = ? AND (title ILIKE ? OR content ILIKE ?)", "published", like, like)
	db.Model(&domain.Poem{}).Count(&total)
	err := db.Preload("Author").Preload("EmotionTags").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, total, err
}

func (r *ReadRepository) GetUserPoems(userID, status string, limit, offset int) ([]*domain.Poem, int64, error) {
	if status == "" {
		status = "published"
	}
	var poems []*domain.Poem
	var total int64
	r.db.Model(&domain.Poem{}).Where("author_id = ? AND status = ?", userID, status).Count(&total)
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("author_id = ? AND status = ?", userID, status).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	return poems, total, err
}

func (r *ReadRepository) GetPoemStats(poemID string) (map[string]interface{}, error) {
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
