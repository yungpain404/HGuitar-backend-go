package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/content/application"
)

func ContentRoutes(contentSvc application.ContentService) chi.Router {
	h := NewContentHandler(contentSvc)
	r := chi.NewRouter()
	r.Mount("/", h.Routes())
	return r
}
