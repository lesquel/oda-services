package http

import (
	"encoding/json"
	"net/http"

	"github.com/lesquel/oda-write-api/internal/pkg/respond"

	"github.com/go-chi/chi/v5"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/validator"
	"github.com/lesquel/oda-write-api/internal/middleware"
	"github.com/lesquel/oda-write-api/internal/features/poems/usecase"
)

// PoemHandler handles poem mutation routes.
type PoemHandler struct{ uc usecase.PoemUseCase }

func NewPoemHandler(uc usecase.PoemUseCase) *PoemHandler { return &PoemHandler{uc: uc} }

func (h *PoemHandler) CreatePoem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req domain.CreatePoemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	poem, err := h.uc.CreatePoem(userID, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, poem)
}

func (h *PoemHandler) UpdatePoem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	poemID := chi.URLParam(r, "id")
	var req domain.UpdatePoemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	poem, err := h.uc.UpdatePoem(poemID, userID, &req)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "unauthorized to update this poem" {
			code = http.StatusForbidden
		}
		respond.Error(w, code, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, poem)
}

func (h *PoemHandler) DeletePoem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	poemID := chi.URLParam(r, "id")
	if err := h.uc.DeletePoem(poemID, userID); err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "unauthorized to delete this poem" {
			code = http.StatusForbidden
		}
		respond.Error(w, code, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PoemHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	poemID := chi.URLParam(r, "id")
	liked, err := h.uc.ToggleLike(poemID, userID)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

func (h *PoemHandler) ToggleBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	poemID := chi.URLParam(r, "id")
	bookmarked, err := h.uc.ToggleBookmark(poemID, userID)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]bool{"bookmarked": bookmarked})
}

func (h *PoemHandler) TagEmotion(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	poemID := chi.URLParam(r, "id")
	var req domain.TagEmotionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.uc.TagEmotion(poemID, userID, req.EmotionID); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "emotion tagged"})
}

func (h *PoemHandler) RemoveEmotionTag(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	poemID := chi.URLParam(r, "id")
	if err := h.uc.RemoveEmotionTag(poemID, userID); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
