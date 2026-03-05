package domain

import (
	"time"

	"gorm.io/gorm"
)

// PoemAuthor contains the author details embedded in a Poem response.
type PoemAuthor struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func (PoemAuthor) TableName() string { return "users" }

// Poem is the core poetry content entity.
type Poem struct {
	ID               string         `gorm:"type:uuid;primaryKey" json:"id"`
	AuthorID         string         `gorm:"type:uuid;not null;index" json:"author_id"`
	Author           *PoemAuthor    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Title            string         `gorm:"not null" json:"title"`
	Content          string         `gorm:"not null" json:"content"`
	Status           string         `gorm:"default:published" json:"status"`
	ModerationStatus string         `gorm:"default:skipped" json:"moderation_status"`
	ModerationScore  float64        `json:"moderation_score,omitempty"`
	ModerationReason string         `json:"moderation_reason,omitempty"`
	ModeratedAt      *time.Time     `json:"moderated_at,omitempty"`
	ModeratedBy      string         `json:"moderated_by,omitempty"`
	LikesCount       int            `gorm:"default:0" json:"likes_count"`
	ViewsCount       int            `gorm:"default:0" json:"views_count"`
	IsLiked          bool           `gorm:"-" json:"is_liked,omitempty"`
	IsBookmarked     bool           `gorm:"-" json:"is_bookmarked,omitempty"`
	UserEmotion      string         `gorm:"-" json:"user_emotion,omitempty"`
	EmotionCounts    map[string]int `gorm:"-" json:"emotion_counts,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	EmotionTags      []EmotionTag   `gorm:"foreignKey:PoemID" json:"emotion_tags,omitempty"`
}

func (Poem) TableName() string { return "poems" }

// Like records that a user liked a poem.
type Like struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	PoemID    string         `gorm:"type:uuid;not null;index" json:"poem_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Like) TableName() string { return "likes" }

// EmotionTag associates a catalog emotion with a poem (by a specific user).
type EmotionTag struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	PoemID    string         `gorm:"type:uuid;not null;index" json:"poem_id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	EmotionID string         `gorm:"type:uuid;not null" json:"emotion_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmotionTag) TableName() string { return "emotion_tags" }

// Bookmark stores a user's saved poem.
type Bookmark struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	PoemID    string         `gorm:"type:uuid;not null;index" json:"poem_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Bookmark) TableName() string { return "bookmarks" }

// EmotionCatalog holds the predefined emotion types users can attach to poems.
type EmotionCatalog struct {
	ID          string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Emoji       string         `json:"emoji"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmotionCatalog) TableName() string { return "emotion_catalog" }

// ModerationLog holds the history of all moderation actions on poems.
type ModerationLog struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	PoemID     string    `gorm:"type:uuid;not null;index" json:"poem_id"`
	Status     string    `gorm:"not null" json:"status"`
	Score      float64   `json:"score"`
	Reason     string    `json:"reason"`
	Provider   string    `json:"provider"`
	Model      string    `json:"model"`
	Categories string    `json:"categories"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ModerationLog) TableName() string { return "moderation_logs" }

// -- Repository interfaces -----------------------------------------------------

type PoemRepository interface {
	Create(poem *Poem) error
	FindByID(id string) (*Poem, error)
	Update(poem *Poem) error
	Delete(id string) error
	GetFeed(limit, offset int) ([]*Poem, error)
	GetUserPoems(userID string, limit, offset int) ([]*Poem, error)
	Search(query string, limit, offset int) ([]*Poem, error)
	IncrementViews(id string) error
	GetStats(poemID string) (map[string]interface{}, error)
}

type LikeRepository interface {
	Toggle(userID, poemID string) (bool, error)
	IsLiked(userID, poemID string) (bool, error)
	GetUserLikes(userID string, limit, offset int) ([]*Poem, error)
}

type EmotionRepository interface {
	Tag(userID, poemID, emotionID string) error
	Remove(userID, poemID, emotionTagID string) error
	GetByPoem(poemID string) ([]*EmotionTag, error)
	GetDistribution(poemID string) (map[string]int, error)
}

type BookmarkRepository interface {
	Toggle(userID, poemID string) (bool, error)
	IsBookmarked(userID, poemID string) (bool, error)
	GetUserBookmarks(userID string, limit, offset int) ([]*Poem, error)
}

type EmotionCatalogRepository interface {
	FindAll() ([]*EmotionCatalog, error)
	FindByID(id string) (*EmotionCatalog, error)
	FindByName(name string) (*EmotionCatalog, error)
	Create(e *EmotionCatalog) error
	Update(e *EmotionCatalog) error
	Delete(id string) error
}

// -- Request / Response DTOs --------------------------------------------------

type CreatePoemRequest struct {
	Title   string `json:"title"   validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
	Status  string `json:"status,omitempty" validate:"omitempty,oneof=published draft"`
}

type UpdatePoemRequest struct {
	Title   string `json:"title"   validate:"omitempty,min=1,max=200"`
	Content string `json:"content" validate:"omitempty,min=1"`
	Status  string `json:"status"  validate:"omitempty,oneof=published draft"`
}

type TagEmotionRequest struct {
	EmotionID string `json:"emotion_id" validate:"required,uuid"`
}

type PoemFeedResponse struct {
	Poems      []*Poem `json:"poems"`
	TotalCount int64   `json:"total_count"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}
