package domain

import "time"

// DashboardStats holds aggregate counts for the admin dashboard.
type DashboardStats struct {
	TotalUsers     int64 `json:"total_users"`
	TotalPoems     int64 `json:"total_poems"`
	TotalLikes     int64 `json:"total_likes"`
	TotalBookmarks int64 `json:"total_bookmarks"`
	TotalEmotions  int64 `json:"total_emotions"`
}

// AdminUser is a user record as seen from the admin panel.
type AdminUser struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Bio       string     `json:"bio"`
	AvatarURL string     `json:"avatar_url"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AdminPoem is a poem record as seen from the admin panel.
type AdminPoem struct {
	ID               string     `json:"id"`
	AuthorID         string     `json:"author_id"`
	Author           *AdminUser `json:"author,omitempty"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	ModerationStatus string     `json:"moderation_status"`
	ModerationScore  float64    `json:"moderation_score,omitempty"`
	ModerationReason string     `json:"moderation_reason,omitempty"`
	ModeratedAt      *time.Time `json:"moderated_at,omitempty"`
	ModeratedBy      string     `json:"moderated_by,omitempty"`
	LikesCount       int        `json:"likes_count"`
	ViewsCount       int        `json:"views_count"`
	CreatedAt        time.Time  `json:"created_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// AdminLike is a like record as seen from the admin panel.
type AdminLike struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PoemID    string     `json:"poem_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AdminBookmark is a bookmark record as seen from the admin panel.
type AdminBookmark struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PoemID    string     `json:"poem_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AdminEmotion is an emotion-tag record as seen from the admin panel.
type AdminEmotion struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PoemID    string     `json:"poem_id"`
	EmotionID string     `json:"emotion_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// PaginatedResponse is a generic paginated list response.
type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}

// -- Admin request DTOs -------------------------------------------------------

type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}

type ChangeStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=published draft removed pending_review rejected"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role"     validate:"required,oneof=user admin"`
}

type UpdateUserAdminRequest struct {
	Username  string `json:"username"   validate:"omitempty,min=3,max=30"`
	Email     string `json:"email"      validate:"omitempty,email"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	IsActive  *bool  `json:"is_active"`
}

type UpdatePoemAdminRequest struct {
	Title   string `json:"title"   validate:"omitempty,min=1,max=200"`
	Content string `json:"content" validate:"omitempty,min=1"`
	Status  string `json:"status"  validate:"omitempty,oneof=published draft removed pending_review rejected"`
}

type CreateEmotionCatalogRequest struct {
	Name        string `json:"name"        validate:"required,min=1,max=50"`
	Emoji       string `json:"emoji"`
	Description string `json:"description"`
}

type UpdateEmotionCatalogRequest struct {
	Name        string `json:"name"        validate:"omitempty,min=1,max=50"`
	Emoji       string `json:"emoji"`
	Description string `json:"description"`
}

type CreateLikeRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	PoemID string `json:"poem_id" validate:"required,uuid"`
}

type CreateBookmarkRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	PoemID string `json:"poem_id" validate:"required,uuid"`
}

type CreateEmotionTagRequest struct {
	UserID    string `json:"user_id"    validate:"required,uuid"`
	PoemID    string `json:"poem_id"    validate:"required,uuid"`
	EmotionID string `json:"emotion_id" validate:"required,uuid"`
}

// AdminModerationLog is a moderation log record as seen from the admin panel.
type AdminModerationLog struct {
	ID         string    `json:"id"`
	PoemID     string    `json:"poem_id"`
	Status     string    `json:"status"`
	Score      float64   `json:"score"`
	Reason     string    `json:"reason"`
	Provider   string    `json:"provider"`
	Model      string    `json:"model"`
	Categories string    `json:"categories"`
	CreatedAt  time.Time `json:"created_at"`
}

type ModerationActionRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Reason string `json:"reason"`
}
