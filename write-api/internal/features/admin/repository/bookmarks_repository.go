package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error) {
	var bms []domain.Bookmark
	db := r.db.Unscoped()
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.Bookmark{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&bms).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminBookmark, len(bms))
	for i, b := range bms {
		items[i] = toAdminBookmark(b)
	}
	return &domain.PaginatedResponse[domain.AdminBookmark]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetBookmark(id string) (*domain.AdminBookmark, error) {
	var b domain.Bookmark
	if err := r.db.Unscoped().First(&b, "id = ?", id).Error; err != nil {
		return nil, err
	}
	ab := toAdminBookmark(b)
	return &ab, nil
}

func (r *adminRepo) CreateBookmark(req *domain.CreateBookmarkRequest) error {
	b := &domain.Bookmark{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		PoemID: req.PoemID,
	}
	return r.db.Create(b).Error
}

func (r *adminRepo) SoftDeleteBookmark(id string) error {
	return r.db.Delete(&domain.Bookmark{}, "id = ?", id).Error
}

func (r *adminRepo) RestoreBookmark(id string) error {
	return r.db.Unscoped().Model(&domain.Bookmark{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *adminRepo) PermanentDeleteBookmark(id string) error {
	return r.db.Unscoped().Delete(&domain.Bookmark{}, "id = ?", id).Error
}
