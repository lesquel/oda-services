package repository

import "github.com/lesquel/oda-shared/domain"

func (r *adminRepo) ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error) {
	var tags []domain.EmotionTag
	db := r.db
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.EmotionTag{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&tags).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminEmotion, len(tags))
	for i, t := range tags {
		items[i] = domain.AdminEmotion{ID: t.ID, UserID: t.UserID, PoemID: t.PoemID, EmotionID: t.EmotionID, CreatedAt: t.CreatedAt}
	}
	return &domain.PaginatedResponse[domain.AdminEmotion]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) HardDeleteEmotion(id string) error {
	return r.db.Delete(&domain.EmotionTag{}, "id = ?", id).Error
}
