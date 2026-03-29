package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents a navigation category
type Category struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Icon        string     `gorm:"size:50" json:"icon,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	IsPublic    bool       `gorm:"default:false" json:"is_public"`
	ShareToken  string     `gorm:"uniqueIndex;size:255" json:"share_token,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Sites []Site `gorm:"foreignKey:CategoryID" json:"sites,omitempty"`
}

// BeforeCreate hook
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
