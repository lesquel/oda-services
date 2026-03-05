package usecase

import (
	"github.com/lesquel/oda-shared/domain"
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
