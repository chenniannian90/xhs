package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email          string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash   string    `gorm:"size:255" json:"-"`
	Username       string    `gorm:"size:100;not null" json:"username"`
	AvatarURL      string    `gorm:"size:500" json:"avatar_url,omitempty"`
	EmailVerified  bool      `gorm:"default:false" json:"email_verified"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	ThemePreference string   `gorm:"size:20;default:'light'" json:"theme_preference"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`

	// Relations
	OAuthAccounts     []OAuthAccount      `gorm:"foreignKey:UserID" json:"-"`
	Categories        []Category          `gorm:"foreignKey:UserID" json:"-"`
	Sites             []Site              `gorm:"foreignKey:UserID" json:"-"`
	EmailVerifications []EmailVerification `gorm:"foreignKey:UserID" json:"-"`
	PasswordResets    []PasswordReset     `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
