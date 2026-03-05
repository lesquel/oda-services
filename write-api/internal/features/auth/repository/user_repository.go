package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

// ── User repository ───────────────────────────────────────────────────────────

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

// ── Refresh token repository ──────────────────────────────────────────────────

type refreshTokenRepo struct{ db *gorm.DB }

// NewRefreshTokenRepository returns a domain.RefreshTokenRepository backed by PostgreSQL.
func NewRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(token *domain.RefreshToken) error {
	if token.ID == "" {
		token.ID = uuid.NewString()
	}
	return r.db.Create(token).Error
}

func (r *refreshTokenRepo) FindByToken(token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	if err := r.db.First(&rt, "token = ? AND expires_at > ?", token, time.Now()).Error; err != nil {
		return nil, errors.New("token not found or expired")
	}
	return &rt, nil
}

func (r *refreshTokenRepo) DeleteByToken(token string) error {
	return r.db.Delete(&domain.RefreshToken{}, "token = ?", token).Error
}

func (r *refreshTokenRepo) DeleteByUserID(userID string) error {
	return r.db.Delete(&domain.RefreshToken{}, "user_id = ?", userID).Error
}
