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

func (uc *adminUseCase) GetLike(id string) (*domain.AdminLike, error) { return uc.repo.GetLike(id) }
func (uc *adminUseCase) CreateLike(req *domain.CreateLikeRequest) error {
	return uc.repo.CreateLike(req)
}
func (uc *adminUseCase) SoftDeleteLike(id string) error      { return uc.repo.SoftDeleteLike(id) }
func (uc *adminUseCase) RestoreLike(id string) error         { return uc.repo.RestoreLike(id) }
func (uc *adminUseCase) PermanentDeleteLike(id string) error { return uc.repo.PermanentDeleteLike(id) }

func (uc *adminUseCase) ListBookmarks(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminBookmark], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListBookmarks(page, limit, poemID, userID)
}

func (uc *adminUseCase) GetBookmark(id string) (*domain.AdminBookmark, error) {
	return uc.repo.GetBookmark(id)
}
func (uc *adminUseCase) CreateBookmark(req *domain.CreateBookmarkRequest) error {
	return uc.repo.CreateBookmark(req)
}
func (uc *adminUseCase) SoftDeleteBookmark(id string) error { return uc.repo.SoftDeleteBookmark(id) }
func (uc *adminUseCase) RestoreBookmark(id string) error    { return uc.repo.RestoreBookmark(id) }
func (uc *adminUseCase) PermanentDeleteBookmark(id string) error {
	return uc.repo.PermanentDeleteBookmark(id)
}

func (uc *adminUseCase) ListEmotions(page, limit int, poemID, userID string) (*domain.PaginatedResponse[domain.AdminEmotion], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListEmotions(page, limit, poemID, userID)
}

func (uc *adminUseCase) GetEmotion(id string) (*domain.AdminEmotion, error) {
	return uc.repo.GetEmotion(id)
}
func (uc *adminUseCase) CreateEmotion(req *domain.CreateEmotionTagRequest) error {
	return uc.repo.CreateEmotion(req)
}
func (uc *adminUseCase) SoftDeleteEmotion(id string) error { return uc.repo.SoftDeleteEmotion(id) }
func (uc *adminUseCase) RestoreEmotion(id string) error    { return uc.repo.RestoreEmotion(id) }
func (uc *adminUseCase) PermanentDeleteEmotion(id string) error {
	return uc.repo.PermanentDeleteEmotion(id)
}
