package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/progress/application"
)

type ProgressHandler struct {
	progressSvc application.ProgressService
}

func NewProgressHandler(progressSvc application.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressSvc: progressSvc}
}

func (h *ProgressHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
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
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.progressSvc.GetUserProgress(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *ProgressHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
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
	lessonID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid lesson id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.progressSvc.CompleteLesson(uid, lessonID, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *ProgressHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/{id}/complete", h.CompleteLesson)
	return r
}
