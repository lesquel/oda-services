package repository

import (
	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

type userRepo struct{ db *gorm.DB }

// NewUserRepository returns a domain.UserRepository backed by PostgreSQL.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	return r.db.Create(user).Error
}

func (r *userRepo) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepo) Search(query string, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	like := "%" + query + "%"
	err := r.db.
		Where("username ILIKE ? OR email ILIKE ?", like, like).
		Limit(limit).Offset(offset).
		Find(&users).Error
	return users, err
}
