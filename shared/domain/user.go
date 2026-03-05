package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents an application user stored in the users table.
type User struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"default:user" json:"role"`
	Bio          string         `json:"bio"`
	AvatarURL    string         `json:"avatar_url"`
	Website      string         `json:"website"`
	Instagram    string         `json:"instagram"`
	Twitter      string         `json:"twitter"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string { return "users" }

// RefreshToken stores long-lived refresh tokens for users.
type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

// IsValid returns true if the token has not expired.
func (r *RefreshToken) IsValid() bool {
	return time.Now().Before(r.ExpiresAt)
}

// -- Repository interfaces -----------------------------------------------------

// UserRepository is the persistence contract for user entities.
type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	Search(query string, limit, offset int) ([]*User, error)
}

// RefreshTokenRepository is the persistence contract for refresh tokens.
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	FindByToken(token string) (*RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID string) error
}

// -- Request / Response DTOs --------------------------------------------------

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	Username  string `json:"username"   validate:"omitempty,min=3,max=30"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Website   string `json:"website"`
	Instagram string `json:"instagram"`
	Twitter   string `json:"twitter"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"   validate:"omitempty,min=3,max=30"`
	AvatarURL string `json:"avatar_url"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}
