package repository

import (
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

// ── Feed ──────────────────────────────────────────────────────────────────────

func (r *ReadRepository) GetFeed(limit, offset int, viewerID string) ([]*domain.Poem, int64, error) {
	var poems []*domain.Poem
	var total int64
	r.db.Model(&domain.Poem{}).Where("status = ?", "published").Count(&total)
	err := r.db.Preload("Author").Preload("EmotionTags").
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	if err == nil {
		err = r.enrichPoems(poems, viewerID)
	}
	return poems, total, err
}

// ── Single poem ───────────────────────────────────────────────────────────────

func (r *ReadRepository) GetPoem(id string, viewerID string) (*domain.Poem, error) {
	var poem domain.Poem
	err := r.db.Preload("Author").Preload("EmotionTags").
		First(&poem, "id = ? AND status = ?", id, "published").Error
	if err != nil {
		return nil, err
	}
	if err := r.enrichPoems([]*domain.Poem{&poem}, viewerID); err != nil {
		return nil, err
	}
	return &poem, nil
}

func (r *ReadRepository) SearchPoems(query string, limit, offset int, viewerID string) ([]*domain.Poem, int64, error) {
	var poems []*domain.Poem
	var total int64
	like := "%" + query + "%"
	db := r.db.Where("status = ? AND (title ILIKE ? OR content ILIKE ?)", "published", like, like)
	db.Model(&domain.Poem{}).Count(&total)
	err := db.Preload("Author").Preload("EmotionTags").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&poems).Error
	if err == nil {
		err = r.enrichPoems(poems, viewerID)
	}
	return poems, total, err
}

func (r *ReadRepository) GetUserPoems(userID, status string, limit, offset int, viewerID string) ([]*domain.Poem, int64, error) {
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
	if err == nil {
		err = r.enrichPoems(poems, viewerID)
	}
	return poems, total, err
}

func (r *ReadRepository) IncrementViews(poemID string) error {
	return r.db.Model(&domain.Poem{}).
		Where("id = ?", poemID).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
}

func (r *ReadRepository) enrichPoems(poems []*domain.Poem, viewerID string) error {
	if len(poems) == 0 {
		return nil
	}

	poemIDs := make([]string, 0, len(poems))
	for _, poem := range poems {
		poemIDs = append(poemIDs, poem.ID)
	}

	type emotionCountRow struct {
		PoemID      string
		EmotionName string
		Count       int
	}
	var emotionRows []emotionCountRow
	if err := r.db.Table("emotion_tags et").
		Select("et.poem_id, ec.name as emotion_name, COUNT(*) as count").
		Joins("JOIN emotion_catalog ec ON ec.id = et.emotion_id").
		Where("et.poem_id IN ?", poemIDs).
		Group("et.poem_id, ec.name").
		Scan(&emotionRows).Error; err != nil {
		return err
	}

	emotionByPoem := make(map[string]map[string]int)
	for _, row := range emotionRows {
		if _, ok := emotionByPoem[row.PoemID]; !ok {
			emotionByPoem[row.PoemID] = make(map[string]int)
		}
		emotionByPoem[row.PoemID][row.EmotionName] = row.Count
	}

	userEmotionByPoem := make(map[string]string)
	likedPoems := make(map[string]bool)
	bookmarkedPoems := make(map[string]bool)

	if viewerID != "" {
		type userEmotionRow struct {
			PoemID      string
			EmotionName string
		}
		var userEmotionRows []userEmotionRow
		if err := r.db.Table("emotion_tags et").
			Select("et.poem_id, ec.name as emotion_name").
			Joins("JOIN emotion_catalog ec ON ec.id = et.emotion_id").
			Where("et.user_id = ? AND et.poem_id IN ?", viewerID, poemIDs).
			Order("et.created_at DESC").
			Scan(&userEmotionRows).Error; err != nil {
			return err
		}
		for _, row := range userEmotionRows {
			if _, exists := userEmotionByPoem[row.PoemID]; !exists {
				userEmotionByPoem[row.PoemID] = row.EmotionName
			}
		}

		type poemRef struct{ PoemID string }
		var likeRows []poemRef
		if err := r.db.Table("likes").
			Select("poem_id").
			Where("user_id = ? AND poem_id IN ?", viewerID, poemIDs).
			Scan(&likeRows).Error; err != nil {
			return err
		}
		for _, row := range likeRows {
			likedPoems[row.PoemID] = true
		}

		var bookmarkRows []poemRef
		if err := r.db.Table("bookmarks").
			Select("poem_id").
			Where("user_id = ? AND poem_id IN ?", viewerID, poemIDs).
			Scan(&bookmarkRows).Error; err != nil {
			return err
		}
		for _, row := range bookmarkRows {
			bookmarkedPoems[row.PoemID] = true
		}
	}

	for _, poem := range poems {
		if counts, ok := emotionByPoem[poem.ID]; ok {
			poem.EmotionCounts = counts
		}
		poem.UserEmotion = userEmotionByPoem[poem.ID]
		poem.IsLiked = likedPoems[poem.ID]
		poem.IsBookmarked = bookmarkedPoems[poem.ID]
	}

	return nil
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
