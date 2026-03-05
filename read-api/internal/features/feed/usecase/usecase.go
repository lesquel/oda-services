package usecase

import (
	"github.com/lesquel/oda-read-api/internal/features/feed/repository"
)

// ReadUseCase handles all read-side business logic.
type ReadUseCase struct {
	repo *repository.ReadRepository
}

func NewReadUseCase(repo *repository.ReadRepository) *ReadUseCase {
	return &ReadUseCase{repo: repo}
}
