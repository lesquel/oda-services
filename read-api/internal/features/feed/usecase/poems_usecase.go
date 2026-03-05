package usecase

import "github.com/lesquel/oda-shared/domain"

func (uc *ReadUseCase) GetFeed(page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetFeed(limit, offset, "")
}

func (uc *ReadUseCase) GetFeedForViewer(page, limit int, viewerID string) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetFeed(limit, offset, viewerID)
}

func (uc *ReadUseCase) GetPoem(id string) (*domain.Poem, error) {
	return uc.GetPoemForViewer(id, "")
}

func (uc *ReadUseCase) GetPoemForViewer(id string, viewerID string) (*domain.Poem, error) {
	poem, err := uc.repo.GetPoem(id, viewerID)
	if err != nil {
		return nil, err
	}
	_ = uc.repo.IncrementViews(id)
	poem.ViewsCount++
	return poem, nil
}

func (uc *ReadUseCase) SearchPoems(query string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.SearchPoems(query, limit, offset, "")
}

func (uc *ReadUseCase) SearchPoemsForViewer(query string, page, limit int, viewerID string) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.SearchPoems(query, limit, offset, viewerID)
}

func (uc *ReadUseCase) GetUserPoems(userID, status string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetUserPoems(userID, status, limit, offset, "")
}

func (uc *ReadUseCase) GetUserPoemsForViewer(userID, status string, page, limit int, viewerID string) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetUserPoems(userID, status, limit, offset, viewerID)
}

func (uc *ReadUseCase) GetPoemStats(poemID string) (map[string]interface{}, error) {
	return uc.repo.GetPoemStats(poemID)
}
