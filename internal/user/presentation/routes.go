package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/user/application"
)

func UserRoutes(userSvc application.UserProfileService) chi.Router {
	h := NewUserHandler(userSvc)
	r := chi.NewRouter()
	r.Mount("/", h.Routes())
	return r
}
