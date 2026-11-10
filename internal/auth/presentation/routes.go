package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/auth/application"
	"github.com/hguitar/backend-go/internal/auth/infrastructure/security"
)

func AuthRoutes(authSvc application.AuthService, firebaseVerifier security.TokenVerifier) chi.Router {
	h := NewAuthHandler(authSvc, firebaseVerifier)
	r := chi.NewRouter()
	r.Mount("/", h.Routes())
	return r
}
