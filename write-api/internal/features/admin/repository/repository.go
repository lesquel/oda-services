package repository

import (
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

// AdminRepository is the full persistence contract for admin operations.
type AdminRepository interface {
	GetDashboardStats() (*domain.DashboardStats, error)

	// Users
	ListUsers(page, limit int, q string) (*domain.PaginatedResponse[domain.AdminUser], error)
	GetUser(id string) (*domain.AdminUser, error)
	CreateUser(req *domain.CreateUserRequest, passwordHash string) error
	UpdateUser(id string, req *domain.UpdateUserAdminRequest) error
	ChangeUserRole(id, role string) error
	SoftDeleteUser(id string) error
	RestoreUser(id string) error
	PermanentDeleteUser(id string) error

	// Poems
	ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error)
	GetPoem(id string) (*domain.AdminPoem, error)
	UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error
	ChangePoemStatus(id, status string) error
	SoftDeletePoem(id string) error
	RestorePoem(id string) error
	PermanentDeletePoem(id string) error

	// Moderation
	ListModerationQueue(page, limit int) (*domain.PaginatedResponse[domain.AdminPoem], error)
	GetModerationLogs(poemID string) ([]domain.AdminModerationLog, error)
	ModerationAction(poemID, action, reason, adminID string) error

	// Likes
	ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error)
	GetLike(id string) (*domain.AdminLike, error)
	CreateLike(req *domain.CreateLikeRequest) error
	SoftDeleteLike(id string) error
	RestoreLike(id string) error
	PermanentDeleteLike(id string) error

	// Bookmarks
	ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error)
	GetBookmark(id string) (*domain.AdminBookmark, error)
	CreateBookmark(req *domain.CreateBookmarkRequest) error
	SoftDeleteBookmark(id string) error
	RestoreBookmark(id string) error
	PermanentDeleteBookmark(id string) error

	// Emotions
	ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error)
	GetEmotion(id string) (*domain.AdminEmotion, error)
	CreateEmotion(req *domain.CreateEmotionTagRequest) error
	SoftDeleteEmotion(id string) error
	RestoreEmotion(id string) error
	PermanentDeleteEmotion(id string) error

	// Emotion Catalog
	ListEmotionCatalog() ([]*domain.EmotionCatalog, error)
	GetEmotionCatalog(id string) (*domain.EmotionCatalog, error)
	CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error
	UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error
	SoftDeleteEmotionCatalog(id string) error
	RestoreEmotionCatalog(id string) error
	PermanentDeleteEmotionCatalog(id string) error
}

type adminRepo struct{ db *gorm.DB }

func NewAdminRepository(db *gorm.DB) AdminRepository { return &adminRepo{db: db} }
