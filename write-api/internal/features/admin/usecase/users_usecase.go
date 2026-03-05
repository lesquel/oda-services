package usecase

import (
	"errors"

	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/hasher"
)

func (uc *adminUseCase) ListUsers(page, limit int, q string) (*domain.PaginatedResponse[domain.AdminUser], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListUsers(page, limit, q)
}

func (uc *adminUseCase) GetUser(id string) (*domain.AdminUser, error) { return uc.repo.GetUser(id) }

func (uc *adminUseCase) CreateUser(req *domain.CreateUserRequest) error {
	hashed, err := hasher.HashPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	return uc.repo.CreateUser(req, hashed)
}

func (uc *adminUseCase) UpdateUser(id string, req *domain.UpdateUserAdminRequest) error {
	return uc.repo.UpdateUser(id, req)
}

func (uc *adminUseCase) ChangeUserRole(id, role string) error {
	return uc.repo.ChangeUserRole(id, role)
}

func (uc *adminUseCase) SoftDeleteUser(id string) error    { return uc.repo.SoftDeleteUser(id) }
func (uc *adminUseCase) RestoreUser(id string) error       { return uc.repo.RestoreUser(id) }
func (uc *adminUseCase) PermanentDeleteUser(id string) error { return uc.repo.PermanentDeleteUser(id) }
