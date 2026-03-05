package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
)

func (r *adminRepo) ListUsers(page, limit int, q string) (*domain.PaginatedResponse[domain.AdminUser], error) {
	var users []domain.User
	db := r.db.Unscoped()
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("username ILIKE ? OR email ILIKE ?", like, like)
	}
	var total int64
	db.Model(&domain.User{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminUser, len(users))
	for i, u := range users {
		items[i] = toAdminUser(u)
	}
	return &domain.PaginatedResponse[domain.AdminUser]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetUser(id string) (*domain.AdminUser, error) {
	var u domain.User
	if err := r.db.Unscoped().First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	au := toAdminUser(u)
	return &au, nil
}

func (r *adminRepo) CreateUser(req *domain.CreateUserRequest, passwordHash string) error {
	u := &domain.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
	}
	return r.db.Create(u).Error
}

func (r *adminRepo) UpdateUser(id string, req *domain.UpdateUserAdminRequest) error {
	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	return r.db.Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) ChangeUserRole(id, role string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("role", role).Error
}

func (r *adminRepo) SoftDeleteUser(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *adminRepo) RestoreUser(id string) error {
	return r.db.Unscoped().Model(&domain.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *adminRepo) PermanentDeleteUser(id string) error {
	return r.db.Unscoped().Delete(&domain.User{}, "id = ?", id).Error
}
