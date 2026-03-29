package repositories

import (
	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/models"
	"gorm.io/gorm"
)

// SiteRepository handles site data operations
type SiteRepository struct {
	db *gorm.DB
}

// NewSiteRepository creates a new site repository
func NewSiteRepository(db *gorm.DB) *SiteRepository {
	return &SiteRepository{db: db}
}

// Create creates a new site
func (r *SiteRepository) Create(site *models.Site) error {
	return r.db.Create(site).Error
}

// FindByID finds a site by ID
func (r *SiteRepository) FindByID(id uuid.UUID) (*models.Site, error) {
	var site models.Site
	err := r.db.Preload("Category").First(&site, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// FindByCategoryID finds all sites in a category
func (r *SiteRepository) FindByCategoryID(categoryID uuid.UUID) ([]models.Site, error) {
	var sites []models.Site
	err := r.db.Where("category_id = ?", categoryID).
		Order("sort_order ASC").
		Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return sites, nil
}

// FindByUserID finds all sites for a user
func (r *SiteRepository) FindByUserID(userID uuid.UUID) ([]models.Site, error) {
	var sites []models.Site
	err := r.db.Where("user_id = ?", userID).
		Preload("Category").
		Order("sort_order ASC").
		Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return sites, nil
}

// Update updates a site
func (r *SiteRepository) Update(site *models.Site) error {
	return r.db.Save(site).Error
}

// Delete deletes a site
func (r *SiteRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Site{}, "id = ?", id).Error
}

// MoveSite moves a site to another category
func (r *SiteRepository) MoveSite(siteID, categoryID uuid.UUID) error {
	return r.db.Model(&models.Site{}).
		Where("id = ?", siteID).
		Update("category_id", categoryID).Error
}

// UpdateSortOrder updates the sort order of a site
func (r *SiteRepository) UpdateSortOrder(id uuid.UUID, sortOrder int) error {
	return r.db.Model(&models.Site{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

// Search searches for sites by name, description, or URL
func (r *SiteRepository) Search(userID uuid.UUID, query string) ([]models.Site, error) {
	var sites []models.Site
	searchQuery := "%" + query + "%"
	err := r.db.Where("user_id = ?", userID).
		Where("name ILIKE ? OR description ILIKE ? OR url ILIKE ?", searchQuery, searchQuery, searchQuery).
		Preload("Category").
		Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return sites, nil
}
