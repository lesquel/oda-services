package usecase

import (
	"errors"

	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/hasher"
)

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

	if req.Username != "" && req.Username != user.Username {
		if existing, _ := uc.userRepo.FindByUsername(req.Username); existing != nil {
			return nil, errors.New("username already taken")
		}
		user.Username = req.Username
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
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
		return errors.New("invalid current password")
	}

	hashed, err := hasher.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = hashed
	if err := uc.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}
	return nil
}

func (uc *authUseCase) GetPublicProfile(username string) (*domain.User, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *authUseCase) SearchUsers(query string) ([]*domain.User, error) {
	return uc.userRepo.Search(query, 20, 0)
}
