package repository

import "github.com/lesquel/oda-shared/domain"

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
