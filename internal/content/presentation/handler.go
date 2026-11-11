package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/content/application"
)

type ContentHandler struct {
	contentSvc application.ContentService
}

func NewContentHandler(contentSvc application.ContentService) *ContentHandler {
	return &ContentHandler{contentSvc: contentSvc}
}

func (h *ContentHandler) GetBlocks(w http.ResponseWriter, r *http.Request) {
	lessonIDStr := chi.URLParam(r, "lessonId")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid lesson id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}
	resp, err := h.contentSvc.GetBlocksForLesson(lessonID, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.LessonBlockResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *ContentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{lessonId}/blocks", h.GetBlocks)
	return r
}
