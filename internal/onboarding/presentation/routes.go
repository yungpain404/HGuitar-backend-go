package presentation

import (
	"github.com/go-chi/chi/v5"
	"github.com/hguitar/backend-go/internal/onboarding/application"
)

func OnboardingRoutes(onboardingSvc application.OnboardingService) chi.Router {
	h := NewOnboardingHandler(onboardingSvc)
	return h.Routes()
}
