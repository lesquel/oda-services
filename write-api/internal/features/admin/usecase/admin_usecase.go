package usecase

import (
	"errors"

	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/hasher"
	"github.com/lesquel/oda-write-api/internal/features/admin/repository"
)

// AdminUseCase wraps all admin operations.
type AdminUseCase interface {
	GetDashboardStats() (*domain.DashboardStats, error)
	ListUsers(page, limit int, q string) (*domain.PaginatedResponse[domain.AdminUser], error)
	GetUser(id string) (*domain.AdminUser, error)
	CreateUser(req *domain.CreateUserRequest) error
	UpdateUser(id string, req *domain.UpdateUserAdminRequest) error
	ChangeUserRole(id, role string) error
	HardDeleteUser(id string) error
	ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error)
	GetPoem(id string) (*domain.AdminPoem, error)
	UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error
	ChangePoemStatus(id, status string) error
	HardDeletePoem(id string) error
	ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error)
	HardDeleteLike(id string) error
	ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error)
	HardDeleteBookmark(id string) error
	ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error)
	HardDeleteEmotion(id string) error
	ListEmotionCatalog() ([]*domain.EmotionCatalog, error)
	CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error
	UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error
	DeleteEmotionCatalog(id string) error
}

type adminUseCase struct{ repo repository.AdminRepository }

func NewAdminUseCase(repo repository.AdminRepository) AdminUseCase {
	return &adminUseCase{repo: repo}
}

func (uc *adminUseCase) GetDashboardStats() (*domain.DashboardStats, error) {
	return uc.repo.GetDashboardStats()
}

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

func (uc *adminUseCase) HardDeleteUser(id string) error { return uc.repo.HardDeleteUser(id) }

func (uc *adminUseCase) ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListPoems(page, limit, q, status)
}

func (uc *adminUseCase) GetPoem(id string) (*domain.AdminPoem, error) { return uc.repo.GetPoem(id) }

func (uc *adminUseCase) UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error {
	return uc.repo.UpdatePoem(id, req)
}

func (uc *adminUseCase) ChangePoemStatus(id, status string) error {
	return uc.repo.ChangePoemStatus(id, status)
}

func (uc *adminUseCase) HardDeletePoem(id string) error { return uc.repo.HardDeletePoem(id) }

func (uc *adminUseCase) ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListLikes(page, limit, poemID, userID)
}

func (uc *adminUseCase) HardDeleteLike(id string) error { return uc.repo.HardDeleteLike(id) }

func (uc *adminUseCase) ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListBookmarks(page, limit, poemID, userID)
}

func (uc *adminUseCase) HardDeleteBookmark(id string) error { return uc.repo.HardDeleteBookmark(id) }

func (uc *adminUseCase) ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListEmotions(page, limit, poemID, userID)
}

func (uc *adminUseCase) HardDeleteEmotion(id string) error { return uc.repo.HardDeleteEmotion(id) }

func (uc *adminUseCase) ListEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	return uc.repo.ListEmotionCatalog()
}

func (uc *adminUseCase) CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error {
	return uc.repo.CreateEmotionCatalog(req)
}

func (uc *adminUseCase) UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error {
	return uc.repo.UpdateEmotionCatalog(id, req)
}

func (uc *adminUseCase) DeleteEmotionCatalog(id string) error {
	return uc.repo.DeleteEmotionCatalog(id)
}
