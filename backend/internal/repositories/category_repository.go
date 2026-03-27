package repositories

import (
	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds a category by ID
func (r *CategoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Sites").First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByUserID finds all categories for a user
func (r *CategoryRepository) FindByUserID(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Sites").Where("user_id = ?", userID).
		Order("sort_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByShareToken finds a public category by share token
func (r *CategoryRepository) FindByShareToken(token string) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Sites").
		Where("share_token = ? AND is_public = ?", token, true).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

// UpdateSortOrder updates the sort order of a category
func (r *CategoryRepository) UpdateSortOrder(id uuid.UUID, sortOrder int) error {
	return r.db.Model(&models.Category{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}
