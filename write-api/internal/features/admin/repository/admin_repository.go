package repository

import (
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"gorm.io/gorm"
)

// AdminRepository is the full persistence contract for admin operations.
type AdminRepository interface {
	GetDashboardStats() (*domain.DashboardStats, error)

	ListUsers(page, limit int, q string) (*domain.PaginatedResponse[domain.AdminUser], error)
	GetUser(id string) (*domain.AdminUser, error)
	CreateUser(req *domain.CreateUserRequest, passwordHash string) error
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

type adminRepo struct{ db *gorm.DB }

func NewAdminRepository(db *gorm.DB) AdminRepository { return &adminRepo{db: db} }

func (r *adminRepo) GetDashboardStats() (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}
	r.db.Model(&domain.User{}).Count(&stats.TotalUsers)
	r.db.Model(&domain.Poem{}).Count(&stats.TotalPoems)
	r.db.Model(&domain.Like{}).Count(&stats.TotalLikes)
	r.db.Model(&domain.Bookmark{}).Count(&stats.TotalBookmarks)
	r.db.Model(&domain.EmotionTag{}).Count(&stats.TotalEmotions)
	return stats, nil
}

// ── Users ────────────────────────────────────────────────────────────────────

func (r *adminRepo) ListUsers(page, limit int, q string) (*domain.PaginatedResponse[domain.AdminUser], error) {
	var users []domain.User
	db := r.db.Unscoped()
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("username ILIKE ? OR email ILIKE ?", like, like)
	}
	var total int64
	db.Model(&domain.User{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminUser, len(users))
	for i, u := range users {
		items[i] = toAdminUser(u)
	}
	return &domain.PaginatedResponse[domain.AdminUser]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetUser(id string) (*domain.AdminUser, error) {
	var u domain.User
	if err := r.db.Unscoped().First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	au := toAdminUser(u)
	return &au, nil
}

func (r *adminRepo) CreateUser(req *domain.CreateUserRequest, passwordHash string) error {
	u := &domain.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
	}
	return r.db.Create(u).Error
}

func (r *adminRepo) UpdateUser(id string, req *domain.UpdateUserAdminRequest) error {
	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	return r.db.Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) ChangeUserRole(id, role string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("role", role).Error
}

func (r *adminRepo) HardDeleteUser(id string) error {
	return r.db.Unscoped().Delete(&domain.User{}, "id = ?", id).Error
}

// ── Poems ────────────────────────────────────────────────────────────────────

func (r *adminRepo) ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	var poems []domain.Poem
	db := r.db.Unscoped().Preload("Author")
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("title ILIKE ? OR content ILIKE ?", like, like)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	db.Model(&domain.Poem{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&poems).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminPoem, len(poems))
	for i, p := range poems {
		items[i] = toAdminPoem(p)
	}
	return &domain.PaginatedResponse[domain.AdminPoem]{
		Items: items, TotalCount: total,
		Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) GetPoem(id string) (*domain.AdminPoem, error) {
	var p domain.Poem
	if err := r.db.Unscoped().Preload("Author").First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	ap := toAdminPoem(p)
	return &ap, nil
}

func (r *adminRepo) UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error {
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) ChangePoemStatus(id, status string) error {
	return r.db.Model(&domain.Poem{}).Where("id = ?", id).Update("status", status).Error
}

func (r *adminRepo) HardDeletePoem(id string) error {
	return r.db.Unscoped().Delete(&domain.Poem{}, "id = ?", id).Error
}

// ── Likes ────────────────────────────────────────────────────────────────────

func (r *adminRepo) ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error) {
	var likes []domain.Like
	db := r.db
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.Like{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&likes).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminLike, len(likes))
	for i, l := range likes {
		items[i] = domain.AdminLike{ID: l.ID, UserID: l.UserID, PoemID: l.PoemID, CreatedAt: l.CreatedAt}
	}
	return &domain.PaginatedResponse[domain.AdminLike]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) HardDeleteLike(id string) error {
	return r.db.Delete(&domain.Like{}, "id = ?", id).Error
}

// ── Bookmarks ────────────────────────────────────────────────────────────────

func (r *adminRepo) ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error) {
	var bms []domain.Bookmark
	db := r.db
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.Bookmark{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&bms).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminBookmark, len(bms))
	for i, b := range bms {
		items[i] = domain.AdminBookmark{ID: b.ID, UserID: b.UserID, PoemID: b.PoemID, CreatedAt: b.CreatedAt}
	}
	return &domain.PaginatedResponse[domain.AdminBookmark]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) HardDeleteBookmark(id string) error {
	return r.db.Delete(&domain.Bookmark{}, "id = ?", id).Error
}

// ── Emotions ─────────────────────────────────────────────────────────────────

func (r *adminRepo) ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error) {
	var tags []domain.EmotionTag
	db := r.db
	if poemID != "" {
		db = db.Where("poem_id = ?", poemID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	var total int64
	db.Model(&domain.EmotionTag{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Find(&tags).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminEmotion, len(tags))
	for i, t := range tags {
		items[i] = domain.AdminEmotion{ID: t.ID, UserID: t.UserID, PoemID: t.PoemID, EmotionID: t.EmotionID, CreatedAt: t.CreatedAt}
	}
	return &domain.PaginatedResponse[domain.AdminEmotion]{
		Items: items, TotalCount: total, Page: page, Limit: limit,
	}, nil
}

func (r *adminRepo) HardDeleteEmotion(id string) error {
	return r.db.Delete(&domain.EmotionTag{}, "id = ?", id).Error
}

// ── Emotion Catalog ──────────────────────────────────────────────────────────

func (r *adminRepo) ListEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	var items []*domain.EmotionCatalog
	return items, r.db.Find(&items).Error
}

func (r *adminRepo) CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error {
	item := &domain.EmotionCatalog{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Emoji:       req.Emoji,
		Description: req.Description,
	}
	return r.db.Create(item).Error
}

func (r *adminRepo) UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error {
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Emoji != "" {
		updates["emoji"] = req.Emoji
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	return r.db.Model(&domain.EmotionCatalog{}).Where("id = ?", id).Updates(updates).Error
}

func (r *adminRepo) DeleteEmotionCatalog(id string) error {
	return r.db.Delete(&domain.EmotionCatalog{}, "id = ?", id).Error
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func toAdminUser(u domain.User) domain.AdminUser {
	au := domain.AdminUser{
		ID: u.ID, Username: u.Username, Email: u.Email,
		Role: u.Role, Bio: u.Bio, AvatarURL: u.AvatarURL,
		IsActive: u.IsActive, CreatedAt: u.CreatedAt,
	}
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		au.DeletedAt = &t
	}
	return au
}

func toAdminPoem(p domain.Poem) domain.AdminPoem {
	ap := domain.AdminPoem{
		ID: p.ID, AuthorID: p.AuthorID, Title: p.Title,
		Content: p.Content, Status: p.Status,
		LikesCount: p.LikesCount, ViewsCount: p.ViewsCount,
		CreatedAt: p.CreatedAt,
	}
	if p.DeletedAt.Valid {
		t := p.DeletedAt.Time
		ap.DeletedAt = &t
	}
	if p.Author != nil {
		ap.Author = &domain.AdminUser{
			ID:        p.Author.ID,
			Username:  p.Author.Username,
			AvatarURL: p.Author.AvatarURL,
		}
	}
	return ap
}

// pageCount is unused but kept for reference.
var (
	_ = math.Ceil
	_ = strings.Contains
)
