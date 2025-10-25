package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"playspotter/internal/models"
	"playspotter/internal/repositories"
	"playspotter/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	tokenRepo *repositories.TokenRepository
	jwtMgr    *jwt.Manager
}

func NewAuthService(userRepo *repositories.UserRepository, tokenRepo *repositories.TokenRepository, jwtMgr *jwt.Manager) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtMgr:    jwtMgr,
	}
}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	// Check if user exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user", // Always create as user
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, errors.New("invalid credentials")
		}
		return "", "", nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", nil, err
	}

	// Generate refresh token
	refreshToken, expiresAt, err := s.jwtMgr.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", nil, err
	}

	// Hash and store refresh token
	tokenHash := hashToken(refreshToken)
	if err := s.tokenRepo.Create(&models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}); err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	// Validate refresh token
	_, err := s.jwtMgr.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Check if token exists and is not revoked
	tokenHash := hashToken(refreshToken)
	storedToken, err := s.tokenRepo.FindByHash(tokenHash)
	if err != nil {
		return "", "", errors.New("refresh token not found or expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(storedToken.UserID)
	if err != nil {
		return "", "", err
	}

	// Revoke old refresh token
	if err := s.tokenRepo.Revoke(tokenHash); err != nil {
		return "", "", err
	}

	// Generate new tokens
	newAccessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, expiresAt, err := s.jwtMgr.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	// Store new refresh token
	newTokenHash := hashToken(newRefreshToken)
	if err := s.tokenRepo.Create(&models.RefreshToken{
		UserID:    user.ID,
		TokenHash: newTokenHash,
		ExpiresAt: expiresAt,
	}); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return s.tokenRepo.Revoke(tokenHash)
}

func (s *AuthService) BootstrapAdmin(email, password string) error {
	// Check if admin already exists
	count, err := s.userRepo.CountByRole("admin")
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("admin already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.User{
		Name:         "Admin",
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	return s.userRepo.Create(admin)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
