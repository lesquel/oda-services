package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/hasher"
	jwtutil "github.com/lesquel/oda-shared/jwt"
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

func (uc *authUseCase) GetProfile(userID string) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *authUseCase) UpdateProfile(userID string, req *domain.UpdateProfileRequest) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	user.Bio = req.Bio
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := uc.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}
	return user, nil
}

func (uc *authUseCase) ChangePassword(userID string, req *domain.ChangePasswordRequest) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !hasher.CheckPassword(user.PasswordHash, req.OldPassword) {
		return errors.New("la contraseña actual es incorrecta")
	}

	newHash, err := hasher.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = newHash
	return uc.userRepo.Update(user)
}

func (uc *authUseCase) GetPublicProfile(username string) (*domain.User, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *authUseCase) SearchUsers(query string) ([]*domain.User, error) {
	users, err := uc.userRepo.Search(query, 20, 0)
	if err != nil {
		return nil, errors.New("failed to search users")
	}
	return users, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

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
