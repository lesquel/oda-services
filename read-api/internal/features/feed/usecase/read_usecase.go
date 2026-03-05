package usecase

import (
	"github.com/lesquel/oda-read-api/internal/features/feed/repository"
	"github.com/lesquel/oda-shared/domain"
)

// ReadUseCase handles all read-side business logic.
type ReadUseCase struct {
	repo *repository.ReadRepository
}

func NewReadUseCase(repo *repository.ReadRepository) *ReadUseCase {
	return &ReadUseCase{repo: repo}
}

func (uc *ReadUseCase) GetFeed(page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetFeed(limit, offset)
}

func (uc *ReadUseCase) GetPoem(id string) (*domain.Poem, error) {
	return uc.repo.GetPoem(id)
}

func (uc *ReadUseCase) SearchPoems(query string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.SearchPoems(query, limit, offset)
}

func (uc *ReadUseCase) GetUserPoems(userID string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetUserPoems(userID, limit, offset)
}

func (uc *ReadUseCase) GetPoemStats(poemID string) (map[string]interface{}, error) {
	return uc.repo.GetPoemStats(poemID)
}

func (uc *ReadUseCase) GetPublicProfile(username string) (*domain.User, error) {
	return uc.repo.GetUserByUsername(username)
}

func (uc *ReadUseCase) SearchUsers(query string, limit, offset int) ([]*domain.User, error) {
	return uc.repo.SearchUsers(query, limit, offset)
}

func (uc *ReadUseCase) GetUserBookmarks(userID string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetUserBookmarks(userID, limit, offset)
}

func (uc *ReadUseCase) GetEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	return uc.repo.GetEmotionCatalog()
}

func (uc *ReadUseCase) GetEmotionDistribution(poemID string) (map[string]int, error) {
	return uc.repo.GetEmotionDistribution(poemID)
}

func (uc *ReadUseCase) GetUserStats(userID string) (map[string]interface{}, error) {
	return uc.repo.GetUserStats(userID)
}
