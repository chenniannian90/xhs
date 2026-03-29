package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("OAuthAccounts").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// UpdateLastLogin updates the last login time
func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", &now).Error
}

// CreateOAuthAccount creates an OAuth account
func (r *UserRepository) CreateOAuthAccount(account *models.OAuthAccount) error {
	return r.db.Create(account).Error
}

// FindOAuthAccount finds an OAuth account by provider and provider user ID
func (r *UserRepository) FindOAuthAccount(provider, providerUserID string) (*models.OAuthAccount, error) {
	var account models.OAuthAccount
	err := r.db.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		Preload("User").
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
