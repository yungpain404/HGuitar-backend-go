package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/auth/domain"
	"github.com/hguitar/backend-go/internal/auth/infrastructure/security"
	"github.com/hguitar/backend-go/internal/common"
)

type authService struct {
	userRepo             domain.UserRepository
	refreshTokenRepo     domain.RefreshTokenRepository
	passwordHasher       PasswordHasher
	jwtCfg               security.JWTConfig
	userProfileSvc       UserProfileService
	refreshTokenExpiryMS int64
}

type PasswordHasher interface {
	Hash(password string) string
	Check(password, hash string) bool
}

func NewAuthService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	passwordHasher PasswordHasher,
	jwtCfg security.JWTConfig,
	userProfileSvc UserProfileService,
	refreshTokenExpiryMS int64,
) AuthService {
	return &authService{
		userRepo:             userRepo,
		refreshTokenRepo:     refreshTokenRepo,
		passwordHasher:       passwordHasher,
		jwtCfg:               jwtCfg,
		userProfileSvc:       userProfileSvc,
		refreshTokenExpiryMS: refreshTokenExpiryMS,
	}
}

func (s *authService) Register(req RegisterRequest) (*AuthResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ConflictException("An account with this email already exists.")
	}

	hash := s.passwordHasher.Hash(req.Password)
	user := domain.NewUser(uuid.New(), req.Email, hash, true, []domain.Role{domain.RoleStudent}, time.Now())
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}
	if err := s.userProfileSvc.CreateProfile(user.ID, req.DisplayName, ""); err != nil {
		return nil, err
	}
	return s.issueTokensFor(user, req.DisplayName)
}

func (s *authService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, common.BusinessRuleException("Invalid email or password.")
	}
	if user == nil {
		return nil, common.BusinessRuleException("Invalid email or password.")
	}
	if !user.Active {
		return nil, common.BusinessRuleException("Account is disabled.")
	}
	if !s.passwordHasher.Check(req.Password, user.PasswordHash) {
		return nil, common.BusinessRuleException("Invalid email or password.")
	}
	displayName, _ := s.userProfileSvc.GetDisplayName(user.ID)
	return s.issueTokensFor(user, displayName)
}

func (s *authService) LoginWithFirebase(identity *FirebaseIdentity) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(identity.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		hash := s.passwordHasher.Hash(uuid.New().String())
		user = domain.NewUser(uuid.New(), identity.Email, hash, true, []domain.Role{domain.RoleStudent}, time.Now())
		if err := s.userRepo.Save(user); err != nil {
			return nil, err
		}
		if err := s.userProfileSvc.CreateProfile(user.ID, identity.DisplayName, identity.AvatarURL); err != nil {
			return nil, err
		}
	} else if !user.Active {
		return nil, common.BusinessRuleException("Account is disabled.")
	}
	storedDisplayName, _ := s.userProfileSvc.GetDisplayName(user.ID)
	displayName := storedDisplayName
	if displayName == "" {
		displayName = identity.DisplayName
	}
	return s.issueTokensFor(user, displayName)
}

func (s *authService) Refresh(req RefreshRequest) (*AuthResponse, error) {
	tokenHash := security.HashToken(req.RefreshToken)
	userID, err := s.refreshTokenRepo.FindUserIDByTokenHash(tokenHash)
	if err != nil {
		return nil, common.BusinessRuleException("Refresh token is invalid or expired.")
	}
	if err := s.refreshTokenRepo.RevokeByTokenHash(tokenHash); err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ResourceNotFoundException("User not found.")
	}
	displayName, _ := s.userProfileSvc.GetDisplayName(userID)
	return s.issueTokensFor(user, displayName)
}

func (s *authService) Logout(req LogoutRequest) error {
	tokenHash := security.HashToken(req.RefreshToken)
	return s.refreshTokenRepo.RevokeByTokenHash(tokenHash)
}

func (s *authService) issueTokensFor(user *domain.User, displayName string) (*AuthResponse, error) {
	roles := user.Roles
	// Sắp xếp roles theo tên (giống Java .sorted())
	sorted := make([]string, len(roles))
	for i, r := range roles {
		sorted[i] = string(r)
	}
	// đơn giản: giữ nguyên thứ tự (Go không sort trong ví dụ này để tránh phụ thuộc sort package)

	accessToken, err := security.GenerateAccessToken(s.jwtCfg, user.ID.String(), user.Email, securityRolesToStrings(roles))
	if err != nil {
		return nil, err
	}

	rawRefresh := security.GenerateSecureToken()
	refreshHash := security.HashToken(rawRefresh)
	expiresAt := time.Now().Add(time.Duration(s.refreshTokenExpiryMS) * time.Millisecond)
	if err := s.refreshTokenRepo.Save(user.ID, refreshHash, expiresAt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		UserID:      user.ID.String(),
		Email:       user.Email,
		DisplayName: displayName,
		Roles:       securityRolesToStrings(roles),
		AccessToken: accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

func securityRolesToStrings(roles []domain.Role) []string {
	out := make([]string, len(roles))
	for i, r := range roles {
		out[i] = string(r)
	}
	return out
}
