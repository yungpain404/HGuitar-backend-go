package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/onboarding/application"
)

type OnboardingHandler struct {
	onboardingSvc application.OnboardingService
	validate      *validator.Validate
}

func NewOnboardingHandler(onboardingSvc application.OnboardingService) *OnboardingHandler {
	return &OnboardingHandler{onboardingSvc: onboardingSvc, validate: validator.New()}
}

func (h *OnboardingHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromCtx(r.Context())
	if userID == "" {
		common.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}
	resp, err := h.onboardingSvc.GetStatus(uid)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *OnboardingHandler) Complete(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromCtx(r.Context())
	if userID == "" {
		common.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}
	var req application.CompleteOnboardingRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	resp, err := h.onboardingSvc.Complete(uid, req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *OnboardingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/status", h.GetStatus)
	r.Post("/complete", h.Complete)
	return r
}
