package usecase

import "github.com/lesquel/oda-shared/domain"

func (uc *adminUseCase) ListEmotionCatalog() ([]*domain.EmotionCatalog, error) {
	return uc.repo.ListEmotionCatalog()
}

func (uc *adminUseCase) GetEmotionCatalog(id string) (*domain.EmotionCatalog, error) {
	return uc.repo.GetEmotionCatalog(id)
}

func (uc *adminUseCase) CreateEmotionCatalog(req *domain.CreateEmotionCatalogRequest) error {
	return uc.repo.CreateEmotionCatalog(req)
}

func (uc *adminUseCase) UpdateEmotionCatalog(id string, req *domain.UpdateEmotionCatalogRequest) error {
	return uc.repo.UpdateEmotionCatalog(id, req)
}

func (uc *adminUseCase) SoftDeleteEmotionCatalog(id string) error {
	return uc.repo.SoftDeleteEmotionCatalog(id)
}

func (uc *adminUseCase) RestoreEmotionCatalog(id string) error {
	return uc.repo.RestoreEmotionCatalog(id)
}

func (uc *adminUseCase) PermanentDeleteEmotionCatalog(id string) error {
	return uc.repo.PermanentDeleteEmotionCatalog(id)
}
