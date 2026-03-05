package usecase

import (
	"errors"

	"github.com/lesquel/oda-shared/domain"
)

func (uc *poemUseCase) CreatePoem(authorID string, req *domain.CreatePoemRequest) (*domain.Poem, error) {
	status := "published"
	if req.Status == "draft" || req.Status == "published" {
		status = req.Status
	}

	poem := &domain.Poem{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
		Status:   status,
	}
	if err := uc.poemRepo.Create(poem); err != nil {
		return nil, errors.New("failed to create poem")
	}
	return poem, nil
}

func (uc *poemUseCase) GetPoemByID(poemID string, userID string) (*domain.Poem, error) {
	poem, err := uc.poemRepo.FindByID(poemID)
	if err != nil {
		return nil, err
	}
	go uc.poemRepo.IncrementViews(poemID) //nolint:errcheck
	return poem, nil
}

func (uc *poemUseCase) GetFeed(limit, offset int, userID string) ([]*domain.Poem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	poems, err := uc.poemRepo.GetFeed(limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch feed")
	}
	return poems, nil
}

func (uc *poemUseCase) GetUserPoems(authorID string, limit, offset int) ([]*domain.Poem, error) {
	poems, err := uc.poemRepo.GetUserPoems(authorID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch user poems")
	}
	return poems, nil
}

func (uc *poemUseCase) UpdatePoem(poemID string, authorID string, req *domain.UpdatePoemRequest) (*domain.Poem, error) {
	poem, err := uc.poemRepo.FindByID(poemID)
	if err != nil {
		return nil, err
	}
	if poem.AuthorID != authorID {
		return nil, errors.New("unauthorized to update this poem")
	}
	if req.Title != "" {
		poem.Title = req.Title
	}
	if req.Content != "" {
		poem.Content = req.Content
	}
	if req.Status != "" {
		poem.Status = req.Status
	}
	if err := uc.poemRepo.Update(poem); err != nil {
		return nil, errors.New("failed to update poem")
	}
	return poem, nil
}

func (uc *poemUseCase) DeletePoem(poemID string, authorID string) error {
	poem, err := uc.poemRepo.FindByID(poemID)
	if err != nil {
		return err
	}
	if poem.AuthorID != authorID {
		return errors.New("unauthorized to delete this poem")
	}
	return uc.poemRepo.Delete(poemID)
}

func (uc *poemUseCase) SearchPoems(query string, limit, offset int) ([]*domain.Poem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return uc.poemRepo.Search(query, limit, offset)
}
