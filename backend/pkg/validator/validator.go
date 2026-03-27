package validator

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/uuid"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error returns the error message
func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// Validator provides validation functions
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// Error returns the validation error message
func (v *Validator) Error() error {
	if len(v.errors) == 0 {
		return nil
	}
	return v.errors[0]
}

// ValidateURL validates a URL
func (v *Validator) ValidateURL(field, value string) {
	if value == "" {
		v.AddError(field, "URL is required")
		return
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		v.AddError(field, "Invalid URL format")
		return
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		v.AddError(field, "URL must include scheme and host")
		return
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		v.AddError(field, "URL must use HTTP or HTTPS scheme")
	}
}

// ValidateName validates a name field
func (v *Validator) ValidateName(field, value string, minLength, maxLength int) {
	if value == "" {
		v.AddError(field, "Name is required")
		return
	}

	if len(value) < minLength {
		v.AddError(field, fmt.Sprintf("Name must be at least %d characters", minLength))
		return
	}

	if len(value) > maxLength {
		v.AddError(field, fmt.Sprintf("Name must not exceed %d characters", maxLength))
		return
	}

	// Check for valid characters (letters, numbers, spaces, common punctuation)
	matched, _ := regexp.MatchString(`^[\p{L}\p{N}\s\-\._\(\)\[\]\{\}]+$`, value)
	if !matched {
		v.AddError(field, "Name contains invalid characters")
	}
}

// ValidateDescription validates a description field
func (v *Validator) ValidateDescription(field, value string, maxLength int) {
	if len(value) > maxLength {
		v.AddError(field, fmt.Sprintf("Description must not exceed %d characters", maxLength))
	}
}

// ValidateIcon validates an icon field (emoji or icon name)
func (v *Validator) ValidateIcon(field, value string) {
	if value == "" {
		return // Icon is optional
	}

	// Allow emojis (multi-byte characters)
	for _, r := range value {
		if r > 127 { // Non-ASCII character (likely emoji)
			continue // Valid emoji
		}
	}

	// Allow short icon names (e.g., "home", "user")
	if len(value) <= 50 && regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`).MatchString(value) {
		return // Valid icon name
	}

	// If icon is neither emoji nor valid name
	if len(value) > 100 {
		v.AddError(field, "Icon is too long")
	}
}

// ValidateSortOrder validates sort order
func (v *Validator) ValidateSortOrder(field string, value int) {
	if value < 0 {
		v.AddError(field, "Sort order must be non-negative")
	}
	if value > 10000 {
		v.AddError(field, "Sort order must not exceed 10000")
	}
}

// ValidateID validates a UUID
func (v *Validator) ValidateID(field, value string) uuid.UUID {
	if value == "" {
		v.AddError(field, "ID is required")
		return uuid.Nil
	}

	id, err := uuid.Parse(value)
	if err != nil {
		v.AddError(field, "Invalid ID format")
		return uuid.Nil
	}

	return id
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(field, value string) {
	if value == "" {
		v.AddError(field, "Email is required")
		return
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.AddError(field, "Invalid email format")
	}
}

// ValidateRequired validates that a field is not empty
func (v *Validator) ValidateRequired(field, value string) {
	if value == "" {
		v.AddError(field, "This field is required")
	}
}

// ValidateMaxLength validates maximum length
func (v *Validator) ValidateMaxLength(field, value string, maxLength int) {
	if len(value) > maxLength {
		v.AddError(field, fmt.Sprintf("Must not exceed %d characters", maxLength))
	}
}

// ValidateMinLength validates minimum length
func (v *Validator) ValidateMinLength(field, value string, minLength int) {
	if len(value) < minLength {
		v.AddError(field, fmt.Sprintf("Must be at least %d characters", minLength))
	}
}

// IsUnique checks if a value is unique in a slice
func IsUnique(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return false
		}
	}
	return true
}

// ValidateUniqueNames checks if category/site names are unique within user's data
func ValidateUniqueNames(newName string, existingNames []string) error {
	if !IsUnique(newName, existingNames) {
		return errors.New("name already exists")
	}
	return nil
}
