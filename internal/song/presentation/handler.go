package presentation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/song/application"
)

type SongHandler struct {
	songSvc application.SongService
	validate *validator.Validate
}

func NewSongHandler(songSvc application.SongService) *SongHandler {
	return &SongHandler{songSvc: songSvc, validate: validator.New()}
}

func (h *SongHandler) ListSongs(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	difficulty := r.URL.Query().Get("difficulty")
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "vi"
	}
	resp, err := h.songSvc.ListSongs(source, difficulty, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		resp = []*application.SongResponse{}
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *SongHandler) GetSong(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid id"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "vi"
	}
	resp, err := h.songSvc.GetSong(id, lang)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	if resp == nil {
		common.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Song not found"})
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var req application.CreateSongRequest
	if err := common.DecodeJSON(r, &req); err != nil {
		common.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	resp, err := h.songSvc.CreateSong(req)
	if err != nil {
		common.WriteError(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, resp)
}

func (h *SongHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListSongs)
	r.Get("/{id}", h.GetSong)
	r.Post("/", h.CreateSong)
	return r
}
