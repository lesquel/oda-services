package repository

import "github.com/lesquel/oda-shared/domain"

func (r *adminRepo) ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error) {
	var bms []domain.Bookmark
	db := r.db
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.Bookmark{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&bms).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminBookmark, len(bms))
	for i, b := range bms {
		items[i] = domain.AdminBookmark{ID: b.ID, UserID: b.UserID, PoemID: b.PoemID, CreatedAt: b.CreatedAt}
	}
	return &domain.PaginatedResponse[domain.AdminBookmark]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) HardDeleteBookmark(id string) error {
	return r.db.Delete(&domain.Bookmark{}, "id = ?", id).Error
}
