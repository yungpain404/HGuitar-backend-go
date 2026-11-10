package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/learning/application"
)

type LearningHandler struct {
	querySvc application.LearningQueryService
	methodSvc application.LearningMethodService
}

func NewLearningHandler(querySvc application.LearningQueryService, methodSvc application.LearningMethodService) *LearningHandler {
	return &LearningHandler{querySvc: querySvc, methodSvc: methodSvc}
}

func (h *LearningHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	guitarType := r.URL.Query().Get("guitarType")
	resp, err := h.querySvc.GetCourses(lang, guitarType)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.CourseResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.querySvc.GetCourse(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Course not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) GetModulesByCourse(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.querySvc.GetModulesByCourse(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.ModuleResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) GetModule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.querySvc.GetModule(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Module not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) GetLessonsByModule(w http.ResponseWriter, r *http.Request) {
	moduleID, err := uuid.Parse(chi.URLParam(r, "moduleId"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid module id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.querySvc.GetLessonsByModule(moduleID, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.LessonResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.querySvc.GetLesson(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Lesson not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) ListMethods(w http.ResponseWriter, r *http.Request) {
	difficulty := r.URL.Query().Get("difficulty")
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "vi"
	}
	resp, err := h.methodSvc.ListMethods(difficulty, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.LearningMethodResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) GetMethod(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "vi"
	}
	resp, err := h.methodSvc.GetMethod(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Method not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *LearningHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/courses", h.GetCourses)
	r.Get("/courses/{id}", h.GetCourse)
	r.Get("/courses/{id}/modules", h.GetModulesByCourse)
	r.Get("/modules/{id}", h.GetModule)
	r.Get("/modules/{moduleId}/lessons", h.GetLessonsByModule)
	r.Get("/lessons/{id}", h.GetLesson)
	r.Get("/learning-methods", h.ListMethods)
	r.Get("/learning-methods/{id}", h.GetMethod)
	return r
}
