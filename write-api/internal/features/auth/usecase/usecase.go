package usecase

import (
	"time"

	"github.com/lesquel/oda-shared/domain"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

// AuthUseCase defines all auth/user mutation operations.
type AuthUseCase interface {
	Register(req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req *domain.LoginRequest) (*domain.AuthResponse, error)
	Refresh(rawToken string) (*domain.AuthResponse, error)
	Logout(rawToken string) error
	GetProfile(userID string) (*domain.User, error)
	UpdateProfile(userID string, req *domain.UpdateProfileRequest) (*domain.User, error)
	ChangePassword(userID string, req *domain.ChangePasswordRequest) error
	GetPublicProfile(username string) (*domain.User, error)
	SearchUsers(query string) ([]*domain.User, error)
}

type authUseCase struct {
	userRepo  domain.UserRepository
	rtRepo    domain.RefreshTokenRepository
	jwtSecret string
}

func NewAuthUseCase(userRepo domain.UserRepository, rtRepo domain.RefreshTokenRepository, jwtSecret string) AuthUseCase {
	return &authUseCase{userRepo: userRepo, rtRepo: rtRepo, jwtSecret: jwtSecret}
}
