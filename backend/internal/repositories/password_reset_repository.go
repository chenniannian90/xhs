package repositories

import (
	"errors"
	"time"

	"github.com/yourusername/navhub/internal/models"
	"gorm.io/gorm"
)

// PasswordResetRepository handles password reset data operations
type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

// FindByToken finds a password reset by token
func (r *PasswordResetRepository) FindByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.Where("token = ?", token).First(&reset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reset, nil
}

// FindByUserID finds all password reset tokens for a user
func (r *PasswordResetRepository) FindByUserID(userID string) ([]models.PasswordReset, error) {
	var resets []models.PasswordReset
	err := r.db.Where("user_id = ?", userID).Find(&resets).Error
	return resets, err
}

// MarkAsUsed marks a token as used
func (r *PasswordResetRepository) MarkAsUsed(reset *models.PasswordReset) error {
	now := time.Now()
	return r.db.Model(reset).Update("used_at", &now).Error
}

// Delete deletes a password reset token
func (r *PasswordResetRepository) Delete(reset *models.PasswordReset) error {
	return r.db.Delete(reset).Error
}

// DeleteByUserID deletes all password reset tokens for a user
func (r *PasswordResetRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.PasswordReset{}).Error
}
