package usecase

import "github.com/lesquel/oda-shared/domain"

func (uc *adminUseCase) ListEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	return uc.repo.ListEmotionCatalog()
}

func (uc *adminUseCase) CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error {
	return uc.repo.CreateEmotionCatalog(req)
}

func (uc *adminUseCase) UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error {
	return uc.repo.UpdateEmotionCatalog(id, req)
}

func (uc *adminUseCase) DeleteEmotionCatalog(id string) error {
	return uc.repo.DeleteEmotionCatalog(id)
}
