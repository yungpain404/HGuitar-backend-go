package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/progress/application"
)

func ProgressRoutes(progressSvc application.ProgressService) chi.Router {
	h := NewProgressHandler(progressSvc)
	return h.Routes()
}
