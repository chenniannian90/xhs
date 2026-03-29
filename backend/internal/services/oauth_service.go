package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/github"
	"github.com/yourusername/navhub/internal/config"
	"github.com/yourusername/navhub/internal/models"
	"github.com/yourusername/navhub/internal/repositories"
)

// OAuthUserInfo represents the user info returned by OAuth providers
type OAuthUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	AvatarURL string `json:"avatar_url"`
}

// OAuthService handles OAuth operations
type OAuthService struct {
	userRepo       *repositories.UserRepository
	authService    *AuthService
	cfg            *config.Config
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	userRepo *repositories.UserRepository,
	authService *AuthService,
	cfg *config.Config,
) *OAuthService {
	return &OAuthService{
		userRepo:    userRepo,
		authService: authService,
		cfg:         cfg,
	}
}

// GetGoogleAuthURL generates the Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	conf := &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  s.cfg.OAuthCallbackURL + "/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return conf.AuthCodeURL(state)
}

// GetGitHubAuthURL generates the GitHub OAuth authorization URL
func (s *OAuthService) GetGitHubAuthURL(state string) string {
	conf := &oauth2.Config{
		ClientID:     s.cfg.GitHubClientID,
		ClientSecret: s.cfg.GitHubClientSecret,
		RedirectURL:  s.cfg.OAuthCallbackURL + "/github/callback",
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	return conf.AuthCodeURL(state)
}

// HandleGoogleCallback handles Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(code string) (*AuthResponse, error) {
	conf := &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  s.cfg.OAuthCallbackURL + "/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from Google
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return s.handleOAuthLogin("google", userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture, token)
}

// HandleGitHubCallback handles GitHub OAuth callback
func (s *OAuthService) HandleGitHubCallback(code string) (*AuthResponse, error) {
	conf := &oauth2.Config{
		ClientID:     s.cfg.GitHubClientID,
		ClientSecret: s.cfg.GitHubClientSecret,
		RedirectURL:  s.cfg.OAuthCallbackURL + "/github/callback",
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from GitHub
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Get primary email from GitHub
	resp, err = client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, fmt.Errorf("failed to get user email: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read email response: %w", err)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return nil, fmt.Errorf("failed to parse emails: %w", err)
	}

	// Find primary verified email
	userInfo.Email = userInfo.Name // Fallback to username
	for _, email := range emails {
		if email.Primary && email.Verified {
			userInfo.Email = email.Email
			break
		}
	}

	return s.handleOAuthLogin("github", userInfo.ID, userInfo.Email, userInfo.Name, userInfo.AvatarURL, token)
}

// handleOAuthLogin handles the login logic for OAuth providers
func (s *OAuthService) handleOAuthLogin(provider, providerUserID, email, name, avatar string, token *oauth2.Token) (*AuthResponse, error) {
	// Find existing OAuth account
	oauthAccount, err := s.userRepo.FindOAuthAccount(provider, providerUserID)

	var user *models.User
	if err == nil && oauthAccount != nil {
		// User exists, update OAuth account
		user = &oauthAccount.User
		now := time.Now()
		oauthAccount.AccessToken = token.AccessToken
		oauthAccount.RefreshToken = token.RefreshToken
		if token.Expiry != (time.Time{}) {
			oauthAccount.ExpiresAt = &now
		}
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	} else {
		// Create new user
		user = &models.User{
			Email:          email,
			Username:        generateUsername(name),
			EmailVerified:   true,
			IsActive:        true,
			ThemePreference: "light",
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}

		// Create OAuth account
		newOAuthAccount := &models.OAuthAccount{
			UserID:        user.ID,
			Provider:       provider,
			ProviderUserID:  providerUserID,
			AccessToken:    token.AccessToken,
			RefreshToken:  token.RefreshToken,
		}
		if token.Expiry != (time.Time{}) {
			now := time.Now()
			newOAuthAccount.ExpiresAt = &now
		}

		if err := s.userRepo.CreateOAuthAccount(newOAuthAccount); err != nil {
			return nil, err
		}
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate auth response using AuthService
	return s.authService.GenerateAuthResponse(user)
}

// GenerateAuthResponse generates auth tokens for a user
func (s *AuthService) GenerateAuthResponse(user *models.User) (*AuthResponse, error) {
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

// GenerateState generates a random state parameter for OAuth
func GenerateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// generateUsername generates a unique username from OAuth user info
func generateUsername(name string) string {
	username := name
	if len(username) > 50 {
		username = username[:50]
	}
	return username + uuid.New().String()[:8]
}
