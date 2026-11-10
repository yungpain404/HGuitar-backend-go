package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/song/application"
)

func SongRoutes(songSvc application.SongService) chi.Router {
	h := NewSongHandler(songSvc)
	return h.Routes()
}
