package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/assessment/application"
)

type QuizHandler struct {
	quizSvc application.QuizService
}

func NewQuizHandler(quizSvc application.QuizService) *QuizHandler {
	return &QuizHandler{quizSvc: quizSvc}
}

func (h *QuizHandler) GetQuizForLesson(w http.ResponseWriter, r *http.Request) {
	lessonID, err := uuid.Parse(chi.URLParam(r, "lessonId"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid lesson id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.quizSvc.GetQuizForLesson(lessonID, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Quiz not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
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
	quizID, err := uuid.Parse(chi.URLParam(r, "quizId"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid quiz id"})
		return
	}
	resp, err := h.quizSvc.StartAttempt(quizID, uid)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusCreated, resp)
}

func (h *QuizHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
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
	attemptID, err := uuid.Parse(chi.URLParam(r, "attemptId"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid attempt id"})
		return
	}
	var req application.SubmitAttemptRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	resp, err := h.quizSvc.SubmitAttempt(attemptID, uid, req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Attempt not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) GetAttempt(w http.ResponseWriter, r *http.Request) {
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
	attemptID, err := uuid.Parse(chi.URLParam(r, "attemptId"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid attempt id"})
		return
	}
	resp, err := h.quizSvc.GetAttempt(attemptID, uid)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Attempt not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/lessons/{lessonId}/quiz", h.GetQuizForLesson)
	r.Post("/quizzes/{quizId}/attempts", h.StartAttempt)
	r.Post("/quiz-attempts/{attemptId}/submit", h.SubmitAttempt)
	r.Get("/quiz-attempts/{attemptId}", h.GetAttempt)
	return r
}
