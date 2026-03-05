package http

import (
	"encoding/json"
	"net/http"

	"github.com/lesquel/oda-write-api/internal/pkg/respond"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/validator"
	"github.com/lesquel/oda-write-api/internal/features/admin/usecase"
)

// AdminHandler handles admin routes.
type AdminHandler struct{ uc usecase.AdminUseCase }

func NewAdminHandler(uc usecase.AdminUseCase) *AdminHandler { return &AdminHandler{uc: uc} }

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetDashboardStats()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, stats)
}

// ── Users ────────────────────────────────────────────────────────────────────

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)
	q := r.URL.Query().Get("q")
	result, err := h.uc.ListUsers(page, limit, q)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.uc.GetUser(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, user)
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.uc.CreateUser(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateUserAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.uc.UpdateUser(chi.URLParam(r, "id"), &req); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *AdminHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	var req domain.ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.uc.ChangeUserRole(chi.URLParam(r, "id"), req.Role); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

func (h *AdminHandler) HardDeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.HardDeleteUser(chi.URLParam(r, "id")); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Poems ────────────────────────────────────────────────────────────────────

func (h *AdminHandler) ListPoems(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.ListPoems(queryInt(r, "page", 1), queryInt(r, "limit", 20),
		r.URL.Query().Get("q"), r.URL.Query().Get("status"))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) GetPoem(w http.ResponseWriter, r *http.Request) {
	poem, err := h.uc.GetPoem(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, poem)
}

func (h *AdminHandler) UpdatePoem(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdatePoemAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.uc.UpdatePoem(chi.URLParam(r, "id"), &req); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "poem updated"})
}

func (h *AdminHandler) ChangePoemStatus(w http.ResponseWriter, r *http.Request) {
	var req domain.ChangeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.uc.ChangePoemStatus(chi.URLParam(r, "id"), req.Status); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}

func (h *AdminHandler) HardDeletePoem(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.HardDeletePoem(chi.URLParam(r, "id")); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Associations ─────────────────────────────────────────────────────────────

func (h *AdminHandler) ListLikes(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.ListLikes(queryInt(r, "page", 1), queryInt(r, "limit", 20),
		r.URL.Query().Get("poem_id"), r.URL.Query().Get("user_id"))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) HardDeleteLike(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.HardDeleteLike(chi.URLParam(r, "id")); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.ListBookmarks(queryInt(r, "page", 1), queryInt(r, "limit", 20),
		r.URL.Query().Get("poem_id"), r.URL.Query().Get("user_id"))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) HardDeleteBookmark(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.HardDeleteBookmark(chi.URLParam(r, "id")); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) ListEmotions(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.ListEmotions(queryInt(r, "page", 1), queryInt(r, "limit", 20),
		r.URL.Query().Get("poem_id"), r.URL.Query().Get("user_id"))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) HardDeleteEmotion(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.HardDeleteEmotion(chi.URLParam(r, "id")); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Emotion Catalog ──────────────────────────────────────────────────────────

func (h *AdminHandler) ListEmotionCatalog(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.ListEmotionCatalog()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, items)
}

func (h *AdminHandler) CreateEmotionCatalog(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateEmotionCatalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.uc.CreateEmotionCatalog(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]string{"message": "emotion created"})
}

func (h *AdminHandler) UpdateEmotionCatalog(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateEmotionCatalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.uc.UpdateEmotionCatalog(chi.URLParam(r, "id"), &req); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "emotion updated"})
}

func (h *AdminHandler) DeleteEmotionCatalog(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.DeleteEmotionCatalog(chi.URLParam(r, "id")); err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func queryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
