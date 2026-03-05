package usecase

import "github.com/lesquel/oda-shared/domain"

func (uc *adminUseCase) ListPoems(page, limit int, q, status string) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListPoems(page, limit, q, status)
}

func (uc *adminUseCase) GetPoem(id string) (*domain.AdminPoem, error) { return uc.repo.GetPoem(id) }

func (uc *adminUseCase) UpdatePoem(id string, req *domain.UpdatePoemAdminRequest) error {
	return uc.repo.UpdatePoem(id, req)
}

func (uc *adminUseCase) ChangePoemStatus(id, status string) error {
	return uc.repo.ChangePoemStatus(id, status)
}

func (uc *adminUseCase) HardDeletePoem(id string) error { return uc.repo.HardDeletePoem(id) }
