package http

import "github.com/lesquel/oda-read-api/internal/features/feed/usecase"

// ReadHandler handles all read-side Huma operations.
type ReadHandler struct{ uc *usecase.ReadUseCase }

func NewReadHandler(uc *usecase.ReadUseCase) *ReadHandler {
	return &ReadHandler{uc: uc}
}
