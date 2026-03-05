package usecase

import "github.com/lesquel/oda-shared/domain"

func (uc *adminUseCase) ListLikes(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminLike], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListLikes(page, limit, poemID, userID)
}

func (uc *adminUseCase) HardDeleteLike(id string) error { return uc.repo.HardDeleteLike(id) }

func (uc *adminUseCase) ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListBookmarks(page, limit, poemID, userID)
}

func (uc *adminUseCase) HardDeleteBookmark(id string) error { return uc.repo.HardDeleteBookmark(id) }

func (uc *adminUseCase) ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListEmotions(page, limit, poemID, userID)
}

func (uc *adminUseCase) HardDeleteEmotion(id string) error { return uc.repo.HardDeleteEmotion(id) }
