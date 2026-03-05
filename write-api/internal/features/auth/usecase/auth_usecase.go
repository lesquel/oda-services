package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/hasher"
	jwtutil "github.com/lesquel/oda-shared/jwt"
)

func (uc *authUseCase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	if existing, _ := uc.userRepo.FindByEmail(req.Email); existing != nil {
		return nil, errors.New("email already registered")
	}
	if existing, _ := uc.userRepo.FindByUsername(req.Username); existing != nil {
		return nil, errors.New("username already taken")
	}

	hashedPwd, err := hasher.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPwd,
		Role:         "user",
		IsActive:     true,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil || !hasher.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}
	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) Refresh(rawToken string) (*domain.AuthResponse, error) {
	rt, err := uc.rtRepo.FindByToken(rawToken)
	if err != nil || !rt.IsValid() {
		return nil, errors.New("invalid or expired refresh token")
	}

	if err := uc.rtRepo.DeleteByToken(rawToken); err != nil {
		return nil, errors.New("failed to revoke old token")
	}

	user, err := uc.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return uc.buildAuthResponse(user)
}

func (uc *authUseCase) Logout(rawToken string) error {
	if rawToken == "" {
		return nil
	}
	return uc.rtRepo.DeleteByToken(rawToken)
}

// buildAuthResponse creates access + refresh tokens for a user.
func (uc *authUseCase) buildAuthResponse(user *domain.User) (*domain.AuthResponse, error) {
	accessToken, err := jwtutil.Generate(user.ID, user.Role, uc.jwtSecret, accessTokenTTL)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	rawRefresh, err := jwtutil.GenerateRefresh(user.ID, user.Role, uc.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	rt := &domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     rawRefresh,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}

	if err := uc.rtRepo.Create(rt); err != nil {
		return nil, errors.New("failed to store refresh token")
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		User:         user,
	}, nil
}
