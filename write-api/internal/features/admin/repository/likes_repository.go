package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error) {
	var likes []domain.Like
	db := r.db.Unscoped()
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.Like{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&likes).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminLike, len(likes))
	for i, l := range likes {
		items[i] = toAdminLike(l)
	}
	return &domain.PaginatedResponse[domain.AdminLike]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetLike(id string) (*domain.AdminLike, error) {
	var l domain.Like
	if err := r.db.Unscoped().First(&l, "id = ?", id).Error; err != nil {
		return nil, err
	}
	al := toAdminLike(l)
	return &al, nil
}

func (r *adminRepo) CreateLike(req *domain.CreateLikeRequest) error {
	l := &domain.Like{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		PoemID: req.PoemID,
	}
	return r.db.Create(l).Error
}

func (r *adminRepo) SoftDeleteLike(id string) error {
	return r.db.Delete(&domain.Like{}, "id = ?", id).Error
}

func (r *adminRepo) RestoreLike(id string) error {
	return r.db.Unscoped().Model(&domain.Like{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *adminRepo) PermanentDeleteLike(id string) error {
	return r.db.Unscoped().Delete(&domain.Like{}, "id = ?", id).Error
}
