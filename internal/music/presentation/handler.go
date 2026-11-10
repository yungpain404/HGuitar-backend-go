package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/music/application"
)

type MusicHandler struct {
	chordSvc       application.ChordService
	strummingSvc   application.StrummingPatternService
}

func NewMusicHandler(chordSvc application.ChordService, strummingSvc application.StrummingPatternService) *MusicHandler {
	return &MusicHandler{chordSvc: chordSvc, strummingSvc: strummingSvc}
}

func (h *MusicHandler) GetChords(w http.ResponseWriter, r *http.Request) {
	guitarType := r.URL.Query().Get("guitarType")
	category := r.URL.Query().Get("category")
	resp, err := h.chordSvc.GetChords(guitarType, category)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.ChordResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) GetChord(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	resp, err := h.chordSvc.GetChord(id)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Chord not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	resp, err := h.chordSvc.GetPositions(id)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.ChordPositionResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) ListPatterns(w http.ResponseWriter, r *http.Request) {
	guitarType := r.URL.Query().Get("guitarType")
	difficulty := r.URL.Query().Get("difficulty")
	resp, err := h.strummingSvc.GetPatterns(guitarType, difficulty)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.StrummingPatternResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) GetPatternDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "vi"
	}
	resp, err := h.strummingSvc.GetPatternDetail(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Pattern not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) GetPatternInstructions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "vi"
	}
	resp, err := h.strummingSvc.GetPatternInstructions(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.StrummingInstructionResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) CreatePattern(w http.ResponseWriter, r *http.Request) {
	var req application.CreateStrummingPatternRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	resp, err := h.strummingSvc.CreatePattern(req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *MusicHandler) ChordRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetChords)
	r.Get("/{id}", h.GetChord)
	r.Get("/{id}/positions", h.GetPositions)
	return r
}

func (h *MusicHandler) StrummingRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListPatterns)
	r.Get("/{id}", h.GetPatternDetail)
	r.Get("/{id}/instructions", h.GetPatternInstructions)
	r.Post("/", h.CreatePattern)
	return r
}
