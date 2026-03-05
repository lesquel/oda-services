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

func (uc *adminUseCase) SoftDeletePoem(id string) error      { return uc.repo.SoftDeletePoem(id) }
func (uc *adminUseCase) RestorePoem(id string) error         { return uc.repo.RestorePoem(id) }
func (uc *adminUseCase) PermanentDeletePoem(id string) error { return uc.repo.PermanentDeletePoem(id) }

// -- Moderation ---------------------------------------------------------------

func (uc *adminUseCase) ListModerationQueue(page, limit int) (*domain.PaginatedResponse[domain.AdminPoem], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repo.ListModerationQueue(page, limit)
}

func (uc *adminUseCase) GetModerationLogs(poemID string) ([]domain.AdminModerationLog, error) {
	return uc.repo.GetModerationLogs(poemID)
}

func (uc *adminUseCase) ModerationAction(poemID, action, reason, adminID string) error {
	return uc.repo.ModerationAction(poemID, action, reason, adminID)
}
