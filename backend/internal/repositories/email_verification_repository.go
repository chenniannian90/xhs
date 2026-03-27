package repositories

import (
	"errors"
	"time"

	"github.com/yourusername/navhub/internal/models"
	"gorm.io/gorm"
)

// EmailVerificationRepository handles email verification data operations
type EmailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository creates a new email verification repository
func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

// Create creates a new email verification token
func (r *EmailVerificationRepository) Create(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

// FindByToken finds an email verification by token
func (r *EmailVerificationRepository) FindByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("token = ?", token).First(&verification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &verification, nil
}

// FindByUserID finds all email verification tokens for a user
func (r *EmailVerificationRepository) FindByUserID(userID string) ([]models.EmailVerification, error) {
	var verifications []models.EmailVerification
	err := r.db.Where("user_id = ?", userID).Find(&verifications).Error
	return verifications, err
}

// Delete deletes an email verification token
func (r *EmailVerificationRepository) Delete(verification *models.EmailVerification) error {
	return r.db.Delete(verification).Error
}

// DeleteByUserID deletes all email verification tokens for a user
func (r *EmailVerificationRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.EmailVerification{}).Error
}

// MarkAsVerified marks a token as verified
func (r *EmailVerificationRepository) MarkAsVerified(verification *models.EmailVerification) error {
	now := time.Now()
	return r.db.Model(verification).Update("verified_at", &now).Error
}
