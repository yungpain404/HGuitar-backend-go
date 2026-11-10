package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/music/application"
)

func MusicRoutes(chordSvc application.ChordService, strummingSvc application.StrummingPatternService) chi.Router {
	h := NewMusicHandler(chordSvc, strummingSvc)
	r := chi.NewRouter()
	r.Mount("/chords", h.ChordRoutes())
	r.Mount("/strumming-patterns", h.StrummingRoutes())
	return r
}
