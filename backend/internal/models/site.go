package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Site represents a navigation site
type Site struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CategoryID  uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	Name        string     `gorm:"size:200;not null" json:"name"`
	URL         string     `gorm:"size:1000;not null" json:"url"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Icon        string     `gorm:"size:500" json:"icon,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	User     User     `gorm:"foreignKey:UserID" json:"-"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate hook
func (s *Site) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
