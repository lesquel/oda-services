package repository

import "github.com/lesquel/oda-shared/domain"

func (r *adminRepo) ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	var poems []domain.Poem
	db := r.db.Unscoped().Preload("Author")
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("title ILIKE ? OR content ILIKE ?", like, like)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	db.Model(&domain.Poem{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&poems).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminPoem, len(poems))
	for i, p := range poems {
		items[i] = toAdminPoem(p)
	}
	return &domain.PaginatedResponse[domain.AdminPoem]{
		Items: items, TotalCount: total,
		Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetPoem(id string) (*domain.AdminPoem, error) {
	var p domain.Poem
	if err := r.db.Unscoped().Preload("Author").First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	ap := toAdminPoem(p)
	return &ap, nil
}

func (r *adminRepo) UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error {
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) ChangePoemStatus(id, status string) error {
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).Update("status", status).Error
}

func (r *adminRepo) HardDeletePoem(id string) error {
	return r.db.Unscoped().Delete(&domain.Poem{}, "id = ?", id).Error
}
