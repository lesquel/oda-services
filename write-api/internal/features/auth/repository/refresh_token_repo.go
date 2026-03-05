package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

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
