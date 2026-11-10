package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/user/application"
)

type UserHandler struct {
	userSvc application.UserProfileService
	validate *validator.Validate
}

func NewUserHandler(userSvc application.UserProfileService) *UserHandler {
	return &UserHandler{
		userSvc:  userSvc,
		validate: validator.New(),
	}
}

func (h *UserHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromCtx(r.Context())
	if userID == "" {
		common.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}
	resp, err := h.userSvc.GetProfile(id)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Profile not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := common.UserIDFromCtx(r.Context())
	if userID == "" {
		common.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}
	var req application.UpdateProfileRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	resp, err := h.userSvc.UpdateProfile(id, req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/me", h.GetMyProfile)
	r.Put("/me", h.UpdateMyProfile)
	return r
}
