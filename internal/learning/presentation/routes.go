package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/learning/application"
)

func LearningRoutes(querySvc application.LearningQueryService, methodSvc application.LearningMethodService) chi.Router {
	h := NewLearningHandler(querySvc, methodSvc)
	return h.Routes()
}
