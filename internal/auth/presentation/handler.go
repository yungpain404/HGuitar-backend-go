package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hguitar/backend-go/internal/auth/application"
	"github.com/hguitar/backend-go/internal/auth/infrastructure/security"
	"github.com/hguitar/backend-go/internal/common"
)

type AuthHandler struct {
	authSvc       application.AuthService
	firebaseVerifier security.TokenVerifier
	validate      *validator.Validate
}

func NewAuthHandler(authSvc application.AuthService, firebaseVerifier security.TokenVerifier) *AuthHandler {
	return &AuthHandler{
		authSvc:         authSvc,
		firebaseVerifier: firebaseVerifier,
		validate:        validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req application.RegisterRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	resp, err := h.authSvc.Register(req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req application.LoginRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	resp, err := h.authSvc.Login(req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Google(w http.ResponseWriter, r *http.Request) {
	var req application.FirebaseLoginRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	identity, err := h.firebaseVerifier.Verify(req.IDToken)
	if err != nil {
		common.WriteError(w, common.BusinessRuleException("Invalid Firebase ID token"))
		return
	}
	resp, err := h.authSvc.LoginWithFirebase(identity)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req application.RefreshRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	resp, err := h.authSvc.Refresh(req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req application.LogoutRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.authSvc.Logout(req); err != nil {
		common.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/google", h.Google)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	return r
}
