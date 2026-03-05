package usecase

import (
	"github.com/lesquel/oda-shared/domain"
)

// PoemUseCase defines poem mutation operations.
type PoemUseCase interface {
	CreatePoem(authorID string, req *domain.CreatePoemRequest) (*domain.Poem, error)
	GetPoemByID(poemID string, userID string) (*domain.Poem, error)
	GetFeed(limit, offset int, userID string) ([]*domain.Poem, error)
	GetUserPoems(authorID string, limit, offset int) ([]*domain.Poem, error)
	UpdatePoem(poemID string, authorID string, req *domain.UpdatePoemRequest) (*domain.Poem, error)
	DeletePoem(poemID string, authorID string) error
	ToggleLike(poemID string, userID string) (bool, error)
	TagEmotion(poemID string, userID string, emotionID string) error
	RemoveEmotionTag(poemID string, userID string) error
	SearchPoems(query string, limit, offset int) ([]*domain.Poem, error)
	ToggleBookmark(poemID string, userID string) (bool, error)
	GetUserBookmarks(userID string, limit, offset int) ([]*domain.Poem, error)
}

type poemUseCase struct {
	poemRepo     domain.PoemRepository
	likeRepo     domain.LikeRepository
	emotionRepo  domain.EmotionRepository
	bookmarkRepo domain.BookmarkRepository
}

func NewPoemUseCase(
	poemRepo domain.PoemRepository,
	likeRepo domain.LikeRepository,
	emotionRepo domain.EmotionRepository,
	bookmarkRepo domain.BookmarkRepository,
) PoemUseCase {
	return &poemUseCase{
		poemRepo:     poemRepo,
		likeRepo:     likeRepo,
		emotionRepo:  emotionRepo,
		bookmarkRepo: bookmarkRepo,
	}
}
