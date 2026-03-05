package repository

import (
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

// ReadRepository combines all read operations needed by the read-api.
type ReadRepository struct {
	db *gorm.DB
}

func NewReadRepository(db *gorm.DB) *ReadRepository {
	return &ReadRepository{db: db}
}

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

// ── Poems ─────────────────────────────────────────────────────────────────────

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

// ── Users ─────────────────────────────────────────────────────────────────────

func (r *ReadRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "username = ? AND is_active = true", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *ReadRepository) SearchUsers(query string, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	like := "%" + query + "%"
	err := r.db.Where("(username ILIKE ? OR email ILIKE ?) AND is_active = true", like, like).
		Limit(limit).Offset(offset).
		Find(&users).Error
	return users, err
}

// ── Bookmarks ─────────────────────────────────────────────────────────────────

func (r *ReadRepository) GetUserBookmarks(userID string, limit, offset int) ([]*domain.Poem, int64, error) {
	var bookmarks []*domain.Bookmark
	var total int64
	r.db.Model(&domain.Bookmark{}).Where("user_id = ?", userID).Count(&total)
	if err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}
	poemIDs := make([]string, len(bookmarks))
	for i, b := range bookmarks {
		poemIDs[i] = b.PoemID
	}
	if len(poemIDs) == 0 {
		return []*domain.Poem{}, total, nil
	}
	var poems []*domain.Poem
	err := r.db.Preload("Author").Where("id IN ?", poemIDs).Find(&poems).Error
	return poems, total, err
}

// ── Emotion Catalog ───────────────────────────────────────────────────────────

func (r *ReadRepository) GetEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Find(&items).Error
}

func (r *ReadRepository) GetEmotionDistribution(poemID string) (map[string]int, error) {
	type row struct {
		EmotionID string
		Count     int
	}
	var rows []row
	err := r.db.Model(&domain.EmotionTag{}).
		Select("emotion_id, COUNT(*) as count").
		Where("poem_id = ?", poemID).
		Group("emotion_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]int, len(rows))
	for _, row := range rows {
		result[row.EmotionID] = row.Count
	}
	return result, nil
}

// ── User Stats ────────────────────────────────────────────────────────────────

func (r *ReadRepository) GetUserStats(userID string) (map[string]interface{}, error) {
	// Poem counts by status
	var poemCount, publishedCount, draftCount int64
	r.db.Model(&domain.Poem{}).Where("author_id = ?", userID).Count(&poemCount)
	r.db.Model(&domain.Poem{}).Where("author_id = ? AND status = ?", userID, "published").Count(&publishedCount)
	r.db.Model(&domain.Poem{}).Where("author_id = ? AND status = ?", userID, "draft").Count(&draftCount)

	// Aggregate likes and views from user's poems
	type aggregates struct {
		TotalLikes int64
		TotalViews int64
	}
	var agg aggregates
	r.db.Model(&domain.Poem{}).
		Select("COALESCE(SUM(likes_count), 0) as total_likes, COALESCE(SUM(views_count), 0) as total_views").
		Where("author_id = ?", userID).
		Scan(&agg)

	// Count bookmarks on user's poems
	var totalBookmarks int64
	r.db.Model(&domain.Bookmark{}).
		Where("poem_id IN (?)", r.db.Model(&domain.Poem{}).Select("id").Where("author_id = ?", userID)).
		Count(&totalBookmarks)

	// Emotion distribution across user's poems
	type emotionRow struct {
		EmotionID string
		Count     int
	}
	var emotionRows []emotionRow
	r.db.Model(&domain.EmotionTag{}).
		Select("emotion_id, COUNT(*) as count").
		Where("poem_id IN (?)", r.db.Model(&domain.Poem{}).Select("id").Where("author_id = ?", userID)).
		Group("emotion_id").
		Scan(&emotionRows)

	emotionDist := make(map[string]int, len(emotionRows))
	for _, row := range emotionRows {
		emotionDist[row.EmotionID] = row.Count
	}

	return map[string]interface{}{
		"poem_count":           poemCount,
		"published_count":      publishedCount,
		"draft_count":          draftCount,
		"total_likes":          agg.TotalLikes,
		"total_views":          agg.TotalViews,
		"total_bookmarks":      totalBookmarks,
		"emotion_distribution": emotionDist,
	}, nil
}
