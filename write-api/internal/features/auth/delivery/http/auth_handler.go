package http

import (
	"encoding/json"
	"net/http"

	"github.com/lesquel/oda-write-api/internal/pkg/respond"

	"github.com/go-chi/chi/v5"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/validator"
	"github.com/lesquel/oda-write-api/internal/middleware"
	"github.com/lesquel/oda-write-api/internal/features/auth/usecase"
)

// AuthHandler handles auth and user mutation routes.
type AuthHandler struct{ uc usecase.AuthUseCase }

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler { return &AuthHandler{uc: uc} }

// logoutRequest is a local DTO for the logout endpoint.
type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.uc.Register(&req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.uc.Login(&req)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.uc.Refresh(req.RefreshToken)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	_ = h.uc.Logout(req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.uc.GetProfile(userID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.uc.UpdateProfile(userID, &req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, user)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req domain.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.ValidateStruct(req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.uc.ChangePassword(userID, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"message": "contraseña actualizada"})
}

func (h *AuthHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	profile, err := h.uc.GetPublicProfile(username)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, profile)
}

func (h *AuthHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		respond.JSON(w, http.StatusOK, []*domain.User{})
		return
	}
	profiles, err := h.uc.SearchUsers(q)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, profiles)
}
