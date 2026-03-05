package usecase

import (
	"errors"

	"github.com/lesquel/oda-shared/domain"
)

func (uc *poemUseCase) ToggleLike(poemID string, userID string) (bool, error) {
	if _, err := uc.poemRepo.FindByID(poemID); err != nil {
		return false, errors.New("poem not found")
	}
	return uc.likeRepo.Toggle(userID, poemID)
}

func (uc *poemUseCase) TagEmotion(poemID string, userID string, emotionID string) error {
	if _, err := uc.poemRepo.FindByID(poemID); err != nil {
		return errors.New("poem not found")
	}
	return uc.emotionRepo.Tag(userID, poemID, emotionID)
}

func (uc *poemUseCase) RemoveEmotionTag(poemID string, userID string) error {
	tags, err := uc.emotionRepo.GetByPoem(poemID)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if tag.UserID == userID {
			return uc.emotionRepo.Remove(userID, poemID, tag.ID)
		}
	}
	return nil // no tag found, nothing to remove
}

func (uc *poemUseCase) ToggleBookmark(poemID string, userID string) (bool, error) {
	if _, err := uc.poemRepo.FindByID(poemID); err != nil {
		return false, errors.New("poem not found")
	}
	return uc.bookmarkRepo.Toggle(userID, poemID)
}

func (uc *poemUseCase) GetUserBookmarks(userID string, limit, offset int) ([]*domain.Poem, error) {
	return uc.bookmarkRepo.GetUserBookmarks(userID, limit, offset)
}
