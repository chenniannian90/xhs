package logger

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	// Log is the global logger instance
	Log *logrus.Logger
)

// Init initializes the logger
func Init(environment string) {
	Log = logrus.New()

	// Set output to stdout
	Log.SetOutput(os.Stdout)

	// Set formatter based on environment
	if environment == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
		Log.SetLevel(logrus.InfoLevel)
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
		Log.SetLevel(logrus.DebugLevel)
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to context
		c.Set("RequestID", requestID)
		c.Set("StartTime", time.Now())

		// Add request ID to response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// LogRequestMiddleware logs HTTP requests
func LogRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID, _ := c.Get("RequestID")

		// Get user ID if authenticated
		userID, exists := c.Get("userID")
		if !exists {
			userID = "anonymous"
		}

		// Build log entry
		entry := Log.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     c.Writer.Status(),
			"latency":    latency.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// Log based on status code
		if c.Writer.Status() >= 500 {
			entry.Error("Server error")
		} else if c.Writer.Status() >= 400 {
			entry.Warn("Client error")
		} else {
			entry.Info("Request completed")
		}
	}
}

// ErrorLogMiddleware logs errors from handlers
func ErrorLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were errors
		if len(c.Errors) > 0 {
			requestID, _ := c.Get("RequestID")

			for _, e := range c.Errors {
				Log.WithFields(logrus.Fields{
					"request_id": requestID,
					"type":       e.Type,
					"error":      e.Error(),
				}).Error("Handler error")
			}
		}
	}
}

// LogWithError logs an error with context
func LogWithError(ctx context.Context, err error, msg string, fields logrus.Fields) {
	requestID := ctx.Value("RequestID")
	if requestID != nil {
		fields["request_id"] = requestID
	}

	Log.WithFields(fields).WithError(err).Error(msg)
}

// LogInfo logs an info message with context
func LogInfo(ctx context.Context, msg string, fields logrus.Fields) {
	requestID := ctx.Value("RequestID")
	if requestID != nil {
		fields["request_id"] = requestID
	}

	Log.WithFields(fields).Info(msg)
}

// LogWarn logs a warning message with context
func LogWarn(ctx context.Context, msg string, fields logrus.Fields) {
	requestID := ctx.Value("RequestID")
	if requestID != nil {
		fields["request_id"] = requestID
	}

	Log.WithFields(fields).Warn(msg)
}

// LogDebug logs a debug message with context
func LogDebug(ctx context.Context, msg string, fields logrus.Fields) {
	requestID := ctx.Value("RequestID")
	if requestID != nil {
		fields["request_id"] = requestID
	}

	Log.WithFields(fields).Debug(msg)
}

// RecoverMiddleware recovers from panics and logs them
func RecoverMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("RequestID")

		Log.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      recovered,
			"stack":      c.Errors.String(),
		}).Error("Panic recovered")

		c.JSON(500, gin.H{
			"error": "Internal server error",
			"code":   500,
		})
	})
}
