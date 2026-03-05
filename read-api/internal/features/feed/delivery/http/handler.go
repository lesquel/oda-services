package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lesquel/oda-read-api/internal/features/feed/usecase"
	"github.com/lesquel/oda-read-api/internal/middleware"
	"github.com/lesquel/oda-read-api/internal/pkg/respond"
)

type ReadHandler struct{ uc *usecase.ReadUseCase }

func NewReadHandler(uc *usecase.ReadUseCase) *ReadHandler {
	return &ReadHandler{uc: uc}
}

func (h *ReadHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)
	poems, total, err := h.uc.GetFeed(page, limit)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"poems": poems, "total_count": total, "page": page, "limit": limit,
	})
}

func (h *ReadHandler) GetPoem(w http.ResponseWriter, r *http.Request) {
	poem, err := h.uc.GetPoem(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, "poem not found")
		return
	}
	respond.JSON(w, http.StatusOK, poem)
}

func (h *ReadHandler) SearchPoems(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)
	poems, total, err := h.uc.SearchPoems(q, page, limit)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"poems": poems, "total_count": total, "page": page, "limit": limit,
	})
}

func (h *ReadHandler) GetUserPoems(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)
	poems, total, err := h.uc.GetUserPoems(userID, page, limit)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"poems": poems, "total_count": total, "page": page, "limit": limit,
	})
}

func (h *ReadHandler) GetPoemStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetPoemStats(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, stats)
}

func (h *ReadHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.uc.GetPublicProfile(chi.URLParam(r, "username"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, "user not found")
		return
	}
	respond.JSON(w, http.StatusOK, user)
}

func (h *ReadHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit := queryInt(r, "limit", 20)
	offset := queryInt(r, "offset", 0)
	users, err := h.uc.SearchUsers(q, limit, offset)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, users)
}

func (h *ReadHandler) GetUserBookmarks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	page := queryInt(r, "page", 1)
	limit := queryInt(r, "limit", 20)
	poems, total, err := h.uc.GetUserBookmarks(userID, page, limit)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"poems": poems, "total_count": total, "page": page, "limit": limit,
	})
}

func (h *ReadHandler) GetEmotionCatalog(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.GetEmotionCatalog()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, items)
}

func (h *ReadHandler) GetEmotionDistribution(w http.ResponseWriter, r *http.Request) {
	dist, err := h.uc.GetEmotionDistribution(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, dist)
}

func (h *ReadHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	stats, err := h.uc.GetUserStats(userID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, stats)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func queryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
