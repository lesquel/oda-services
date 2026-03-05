package usecase

import "github.com/lesquel/oda-shared/domain"

func (uc *ReadUseCase) GetPublicProfile(username string) (*domain.User, error) {
	return uc.repo.GetUserByUsername(username)
}

func (uc *ReadUseCase) SearchUsers(query string, limit, offset int) ([]*domain.User, error) {
	return uc.repo.SearchUsers(query, limit, offset)
}

func (uc *ReadUseCase) GetUserStats(userID string) (map[string]interface{}, error) {
	return uc.repo.GetUserStats(userID)
}

func (uc *ReadUseCase) GetUserBookmarks(userID string, page, limit int) ([]*domain.Poem, int64, error) {
	offset := (page - 1) * limit
	return uc.repo.GetUserBookmarks(userID, limit, offset, userID)
}

func (uc *ReadUseCase) GetEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	return uc.repo.GetEmotionCatalog()
}

func (uc *ReadUseCase) GetEmotionDistribution(poemID string) (map[string]int, error) {
	return uc.repo.GetEmotionDistribution(poemID)
}
