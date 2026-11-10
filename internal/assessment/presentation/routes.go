package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/assessment/application"
)

func QuizRoutes(quizSvc application.QuizService) chi.Router {
	h := NewQuizHandler(quizSvc)
	return h.Routes()
}
