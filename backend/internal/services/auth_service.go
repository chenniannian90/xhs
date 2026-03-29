package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/config"
	"github.com/yourusername/navhub/internal/models"
	"github.com/yourusername/navhub/internal/repositories"
	"github.com/yourusername/navhub/pkg/jwt"
	"github.com/yourusername/navhub/pkg/password"
	"github.com/yourusername/navhub/pkg/validator"
)

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *models.User `json:"user"`
}

// AuthService handles authentication business logic
type AuthService struct {
	userRepo                *repositories.UserRepository
	emailVerificationRepo   *repositories.EmailVerificationRepository
	passwordResetRepo       *repositories.PasswordResetRepository
	jwtMgr                  *jwt.Manager
	emailSvc                *EmailService
	jwtExpiry               int
	refreshExpiry           int
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repositories.UserRepository,
	emailVerificationRepo *repositories.EmailVerificationRepository,
	passwordResetRepo *repositories.PasswordResetRepository,
	cfg *config.Config,
	emailSvc *EmailService,
) *AuthService {
	return &AuthService{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		passwordResetRepo:     passwordResetRepo,
		jwtMgr:                jwt.NewManager(cfg.JWTSecret),
		emailSvc:              emailSvc,
		jwtExpiry:             cfg.JWTExpiry,
		refreshExpiry:         cfg.RefreshExpiry,
	}
}

// RegisterInput represents registration input
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Username string `json:"username" binding:"required,min=3,max=50"`
}

// Register registers a new user
func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	// Validate input
	v := validator.NewValidator()
	v.ValidateEmail("email", input.Email)
	v.ValidateName("username", input.Username, 3, 50)
	v.ValidateMinLength("password", input.Password, 8)
	v.ValidateMaxLength("password", input.Password, 128)

	if v.HasErrors() {
		return nil, v.Error()
	}

	// Validate password strength
	if err := password.ValidateStrength(input.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hash, err := password.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:           input.Email,
		PasswordHash:    hash,
		Username:        input.Username,
		EmailVerified:   false,
		IsActive:        true,
		ThemePreference: "light",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate verification token
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	// Create email verification record (expires in 24 hours)
	verification := &models.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.emailVerificationRepo.Create(verification); err != nil {
		return nil, err
	}

	// Send verification email
	if err := s.emailSvc.SendVerificationEmail(user.Email, token); err != nil {
		// Log error but don't fail registration
		_ = err
	}

	// Generate JWT tokens
	accessToken, err := s.jwtMgr.GenerateToken(user.ID, user.Email, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtMgr.GenerateRefreshToken(user.ID, user.Email, s.refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// LoginInput represents login input
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user
func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !password.VerifyPassword(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens
	accessToken, err := s.jwtMgr.GenerateToken(user.ID, user.Email, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtMgr.GenerateRefreshToken(user.ID, user.Email, s.refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	// Parse string to UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}
	return s.userRepo.FindByID(id)
}

// VerifyEmail verifies a user's email
func (s *AuthService) VerifyEmail(token string) error {
	// Find verification token by token
	verification, err := s.emailVerificationRepo.FindByToken(token)
	if err != nil {
		return errors.New("invalid verification token")
	}

	if verification == nil {
		return errors.New("invalid verification token")
	}

	// Check if expired
	if verification.IsExpired() {
		return errors.New("verification token has expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(verification.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update user email_verified = true
	user.EmailVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Mark verification as verified
	if err := s.emailVerificationRepo.MarkAsVerified(verification); err != nil {
		return err
	}

	// Delete verification token
	s.emailVerificationRepo.Delete(verification)

	return nil
}

// ForgotPasswordInput represents forgot password input
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword initiates password reset
func (s *AuthService) ForgotPassword(input ForgotPasswordInput) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	// Generate reset token
	token, err := GenerateToken()
	if err != nil {
		return err
	}

	// Create password reset record (expires in 1 hour)
	reset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.passwordResetRepo.Create(reset); err != nil {
		return err
	}

	// Send reset email
	if err := s.emailSvc.SendPasswordResetEmail(user.Email, token); err != nil {
		// Log error but don't fail the request
		_ = err
	}

	return nil
}

// ResetPasswordInput represents reset password input
type ResetPasswordInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(input ResetPasswordInput) error {
	// Find reset token by token
	reset, err := s.passwordResetRepo.FindByToken(input.Token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	if reset == nil {
		return errors.New("invalid reset token")
	}

	// Check if expired or used
	if reset.IsExpired() {
		return errors.New("reset token has expired")
	}

	if reset.IsUsed() {
		return errors.New("reset token has already been used")
	}

	// Get user
	user, err := s.userRepo.FindByID(reset.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate password strength
	if err := password.ValidateStrength(input.Password); err != nil {
		return err
	}

	// Hash new password
	hash, err := password.HashPassword(input.Password)
	if err != nil {
		return err
	}

	// Update user password
	user.PasswordHash = hash
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(reset); err != nil {
		return err
	}

	return nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return s.jwtMgr.ValidateToken(tokenString)
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	claims, err := s.jwtMgr.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := s.jwtMgr.GenerateToken(user.ID, user.Email, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtMgr.GenerateRefreshToken(user.ID, user.Email, s.refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}
