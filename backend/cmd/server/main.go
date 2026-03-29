package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/navhub/internal/config"
	"github.com/yourusername/navhub/internal/handlers"
	"github.com/yourusername/navhub/internal/middleware"
	"github.com/yourusername/navhub/internal/models"
	"github.com/yourusername/navhub/internal/repositories"
	"github.com/yourusername/navhub/internal/services"
	"github.com/yourusername/navhub/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Environment)
	logger.Log.Info("Starting NavHub Backend API")
	logger.Log.WithFields(map[string]interface{}{
		"environment": cfg.Environment,
		"port":        cfg.Port,
		"database":    cfg.DatabaseName,
	}).Info("Configuration loaded")

	// Initialize security configuration
	securityConfig := middleware.DefaultSecurityConfig(cfg.FrontendURL)
	if cfg.Environment == "production" {
		securityConfig.HSTSEnabled = true
		securityConfig.FrameOptions = "SAMEORIGIN"
		securityConfig.AllowedOrigins = []string{cfg.FrontendURL}
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	fmt.Println("✅ Database connected successfully!")
	fmt.Printf("📦 Database: %s\n", cfg.DatabaseName)
	fmt.Printf("🔌 Port: %d\n", cfg.Port)

	// Auto migrate tables
	fmt.Println("🔄 Running database migrations...")
	for _, model := range models.Models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate %T: %v", model, err)
		}
	}
	fmt.Println("✅ Database migrations completed!")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	siteRepo := repositories.NewSiteRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)

	// Initialize services
	emailService := services.NewEmailService(cfg)
	authService := services.NewAuthService(userRepo, emailVerificationRepo, passwordResetRepo, cfg, emailService)
	oauthService := services.NewOAuthService(userRepo, authService, cfg)
	categoryService := services.NewCategoryService(categoryRepo, siteRepo)
	siteService := services.NewSiteService(siteRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, oauthService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	siteHandler := handlers.NewSiteHandler(siteService)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Apply security middleware
	router.Use(middleware.SecurityHeaders(securityConfig))
	router.Use(middleware.EnhancedCORS(securityConfig))
	router.Use(middleware.RequestSizeLimit(10 << 20)) // 10MB max request size
	router.Use(middleware.StripSlashes())

	// Apply logging middleware
	router.Use(logger.RequestIDMiddleware())
	router.Use(logger.RecoverMiddleware())
	router.Use(logger.LogRequestMiddleware())
	router.Use(logger.ErrorLogMiddleware())

	// Apply rate limiting
	rateLimitConfig := middleware.DefaultRateLimitConfig()
	if cfg.Environment == "production" {
		// Stricter rate limits in production
		rateLimitConfig = middleware.RateLimitConfig{
			RequestsPerMinute: 120, // 120 requests per minute
			Burst:             20,
			SkipRateLimit: func(ip string) bool {
				return ip == "127.0.0.1" || ip == "::1"
			},
		}
	}
	router.Use(middleware.RateLimit(rateLimitConfig))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"database": "connected",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.GET("/oauth/:provider", authHandler.OAuthLogin)
			auth.GET("/oauth/:provider/callback", authHandler.OAuthCallback)

			// Protected auth routes
			authAuth := auth.Group("")
			authAuth.Use(middleware.Auth(authService))
			{
				authAuth.GET("/me", authHandler.GetCurrentUser)
			}
		}

	// Category routes (protected)
		categories := v1.Group("/categories")
		categories.Use(middleware.Auth(authService))
		{
			categories.GET("", categoryHandler.ListCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/search", categoryHandler.SearchCategories)
			categories.GET("/export", categoryHandler.ExportCategories)
			categories.POST("/import", categoryHandler.ImportCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/:id/export", categoryHandler.ExportCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.POST("/:id/share", categoryHandler.ShareCategory)
			categories.DELETE("/:id/share", categoryHandler.UnshareCategory)
			categories.POST("/reorder", categoryHandler.ReorderCategories)
			categories.DELETE("/batch", categoryHandler.BatchDeleteCategories)
			categories.PUT("/batch", categoryHandler.BatchUpdateCategories)
		}

		// Public shared category
		v1.GET("/shared/:token", categoryHandler.GetSharedCategory)

		// Site routes (protected)
		sites := v1.Group("/sites")
		sites.Use(middleware.Auth(authService))
		{
			sites.GET("", siteHandler.ListSites)
			sites.GET("/search", siteHandler.SearchSites)
			sites.POST("", siteHandler.CreateSite)
			sites.GET("/:id", siteHandler.GetSite)
			sites.PUT("/:id", siteHandler.UpdateSite)
			sites.DELETE("/:id", siteHandler.DeleteSite)
			sites.PUT("/:id/move", siteHandler.MoveSite)
			sites.POST("/batch", siteHandler.BatchCreateSites)
			sites.POST("/reorder", siteHandler.ReorderSites)
			sites.DELETE("/batch", siteHandler.BatchDelete)
			sites.PUT("/batch", siteHandler.BatchUpdate)
		}

		// Search routes (protected)
		search := v1.Group("/search")
		search.Use(middleware.Auth(authService))
		{
			search.GET("/sites", siteHandler.GlobalSearchSites)
			search.GET("/categories", categoryHandler.SearchCategories)
		}
	}

	// Server configuration
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Log.WithFields(map[string]interface{}{
			"port": cfg.Port,
			"api":  fmt.Sprintf("http://localhost:%d/api/v1", cfg.Port),
		}).Info("Server starting")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.WithError(err).Error("Failed to start server")
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Log.Info("Waiting for connections to close...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Error("Server forced to shutdown")
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Log.Info("Server exited gracefully")
}
