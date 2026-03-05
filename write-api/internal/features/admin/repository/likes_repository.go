package repository

import "github.com/lesquel/oda-shared/domain"

func (r *adminRepo) ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error) {
	var likes []domain.Like
	db := r.db
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.Like{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&likes).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminLike, len(likes))
	for i, l := range likes {
		items[i] = domain.AdminLike{ID: l.ID, UserID: l.UserID, PoemID: l.PoemID, CreatedAt: l.CreatedAt}
	}
	return &domain.PaginatedResponse[domain.AdminLike]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) HardDeleteLike(id string) error {
	return r.db.Delete(&domain.Like{}, "id = ?", id).Error
}
