package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/navhub/internal/services"
)

// AuthHandler handles auth HTTP requests
type AuthHandler struct {
	authService *services.AuthService
	oauthService *services.OAuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService, oauthService *services.OAuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string   `json:"error"`
	Code        int      `json:"code,omitempty"`
	Details     string   `json:"details,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.authService.Register(input)
	if err != nil {
		// Check if error is about email already registered
		errorMsg := err.Error()
		if errorMsg == "email already registered" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "该邮箱已被注册",
				Details: fmt.Sprintf("邮箱 %s 已经被注册，您可以直接登录", input.Email),
				Suggestions: []string{
					"直接登录：前往登录页面使用您的邮箱和密码登录",
					"忘记密码：如果忘记密码，可以通过忘记密码功能重置",
					"使用其他邮箱：使用不同的邮箱地址创建新账号",
				},
			})
			return
		}

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: errorMsg})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    201,
		Message: "Registration successful. Please check your email to verify your account.",
		Data:    user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	response, err := h.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Login successful",
		Data:    response,
	})
}

// GetCurrentUser gets the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	// Get user from service
	user, err := h.authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    user,
	})
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.VerifyEmail(req.Token); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Email verified successfully",
	})
}

// ForgotPassword handles password reset requests
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input services.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.ForgotPassword(input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Password reset email sent",
	})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input services.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.ResetPassword(input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Password reset successful",
	})
}

// RefreshToken refreshes an access token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Token refreshed",
		Data:    response,
	})
}

// OAuthLogin handles OAuth login request
func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")

	var authURL string

	switch provider {
	case "google":
		state, err := services.GenerateState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate state"})
			return
		}
		authURL = h.oauthService.GetGoogleAuthURL(state)
	case "github":
		state, err := services.GenerateState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate state"})
			return
		}
		authURL = h.oauthService.GetGitHubAuthURL(state)
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid OAuth provider"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "OAuth authorization URL generated",
		Data:    gin.H{"auth_url": authURL},
	})
}

// OAuthCallback handles OAuth callback
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Authorization code is required"})
		return
	}

	var response *services.AuthResponse
	var err error

	switch provider {
	case "google":
		response, err = h.oauthService.HandleGoogleCallback(code)
	case "github":
		response, err = h.oauthService.HandleGitHubCallback(code)
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid OAuth provider"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "OAuth login successful",
		Data:    response,
	})
}
