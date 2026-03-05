package usecase

import "github.com/lesquel/oda-shared/domain"

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

func (uc *ReadUseCase) GetUserPoems(userID, status string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetUserPoems(userID, status, limit, offset)
}

func (uc *ReadUseCase) GetPoemStats(poemID string) (map[string]interface{}, error) {
	return uc.repo.GetPoemStats(poemID)
}
