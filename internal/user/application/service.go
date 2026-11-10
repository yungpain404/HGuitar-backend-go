package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/user/domain"
)

type UserProfileService interface {
	CreateProfile(userID uuid.UUID, displayName, avatarURL string) error
	GetProfile(userID uuid.UUID) (*UserProfileResponse, error)
	UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*UserProfileResponse, error)
	GetDisplayName(userID uuid.UUID) (string, error)
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=1"`
	AvatarURL   string `json:"avatar_url" validate:"omitempty,url"`
	GuitarType  string `json:"guitar_type" validate:"omitempty"`
}

type UserProfileResponse struct {
	UserID    string `json:"user_id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	GuitarType  string `json:"guitar_type"`
}

type userProfileService struct {
	repo domain.UserProfileRepository
}

func NewUserProfileService(repo domain.UserProfileRepository) UserProfileService {
	return &userProfileService{repo: repo}
}

func (s *userProfileService) CreateProfile(userID uuid.UUID, displayName, avatarURL string) error {
	now := time.Now()
	existing, err := s.repo.FindByUserID(userID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil // đã có profile
	}
	profile := domain.NewUserProfile(uuid.New(), userID, displayName, avatarURL, "", now)
	return s.repo.Save(profile)
}

func (s *userProfileService) GetProfile(userID uuid.UUID) (*UserProfileResponse, error) {
	p, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return &UserProfileResponse{
		UserID:       p.UserID.String(),
		DisplayName:  p.DisplayName,
		AvatarURL:    p.AvatarURL,
		GuitarType:   p.GuitarType,
	}, nil
}

func (s *userProfileService) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*UserProfileResponse, error) {
	p, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		p = domain.NewUserProfile(uuid.New(), userID, "", "", "", time.Now())
	}
	if req.DisplayName != "" {
		p.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		p.AvatarURL = req.AvatarURL
	}
	if req.GuitarType != "" {
		p.GuitarType = req.GuitarType
	}
	p.UpdatedAt = time.Now()
	if err := s.repo.Save(p); err != nil {
		return nil, err
	}
	return &UserProfileResponse{
		UserID:       p.UserID.String(),
		DisplayName:  p.DisplayName,
		AvatarURL:    p.AvatarURL,
		GuitarType:   p.GuitarType,
	}, nil
}

func (s *userProfileService) GetDisplayName(userID uuid.UUID) (string, error) {
	p, err := s.repo.FindByUserID(userID)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", nil
	}
	return p.DisplayName, nil
}
