package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/onboarding/domain"
)

type OnboardingStatusResponse struct {
	UserID    string `json:"user_id"`
	Completed bool   `json:"completed"`
	GuitarType string `json:"guitar_type,omitempty"`
}

type CompleteOnboardingRequest struct {
	GuitarType string `json:"guitar_type" validate:"required"`
}

type OnboardingService interface {
	GetStatus(userID uuid.UUID) (*OnboardingStatusResponse, error)
	Complete(userID uuid.UUID, req CompleteOnboardingRequest) (*OnboardingStatusResponse, error)
}

type onboardingService struct {
	repo domain.UserPreferenceRepository
}

func NewOnboardingService(repo domain.UserPreferenceRepository) OnboardingService {
	return &onboardingService{repo: repo}
}

func (s *onboardingService) GetStatus(userID uuid.UUID) (*OnboardingStatusResponse, error) {
	pref, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		return &OnboardingStatusResponse{UserID: userID.String(), Completed: false}, nil
	}
	return &OnboardingStatusResponse{
		UserID:    userID.String(),
		Completed: pref.Completed,
		GuitarType: pref.GuitarType,
	}, nil
}

func (s *onboardingService) Complete(userID uuid.UUID, req CompleteOnboardingRequest) (*OnboardingStatusResponse, error) {
	pref, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if pref == nil {
		pref = &domain.UserPreference{
			ID:   uuid.New(),
			UserID: userID,
			GuitarType: req.GuitarType,
			Completed: true,
			CreatedAt: now,
			UpdatedAt: now,
		}
	} else {
		pref.GuitarType = req.GuitarType
		pref.Completed = true
		pref.UpdatedAt = now
	}
	if err := s.repo.Save(pref); err != nil {
		return nil, err
	}
	return &OnboardingStatusResponse{
		UserID:    userID.String(),
		Completed: true,
		GuitarType: pref.GuitarType,
	}, nil
}
