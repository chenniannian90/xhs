package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/models"
	"github.com/yourusername/navhub/internal/repositories"
	"github.com/yourusername/navhub/pkg/validator"
)

// SiteService handles site business logic
type SiteService struct {
	siteRepo *repositories.SiteRepository
}

// NewSiteService creates a new site service
func NewSiteService(siteRepo *repositories.SiteRepository) *SiteService {
	return &SiteService{
		siteRepo: siteRepo,
	}
}

// CreateSiteInput represents site creation input
type CreateSiteInput struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Name        string    `json:"name" binding:"required,max=200"`
	URL         string    `json:"url" binding:"required,url"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	SortOrder   int       `json:"sort_order"`
}

// Create creates a new site
func (s *SiteService) Create(userID uuid.UUID, input CreateSiteInput) (*models.Site, error) {
	// Validate input
	v := validator.NewValidator()
	v.ValidateName("name", input.Name, 1, 200)
	v.ValidateURL("url", input.URL)
	v.ValidateDescription("description", input.Description, 500)
	v.ValidateIcon("icon", input.Icon)
	v.ValidateSortOrder("sort_order", input.SortOrder)

	if v.HasErrors() {
		return nil, v.Error()
	}

	site := &models.Site{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		URL:         input.URL,
		Description: input.Description,
		Icon:        input.Icon,
		SortOrder:   input.SortOrder,
	}

	if err := s.siteRepo.Create(site); err != nil {
		return nil, err
	}

	return site, nil
}

// GetByID gets a site by ID
func (s *SiteService) GetByID(id, userID uuid.UUID) (*models.Site, error) {
	site, err := s.siteRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if site.UserID != userID {
		return nil, errors.New("access denied")
	}

	return site, nil
}

// ListByUser gets all sites for a user
func (s *SiteService) ListByUser(userID uuid.UUID) ([]models.Site, error) {
	return s.siteRepo.FindByUserID(userID)
}

// ListByCategory gets all sites in a category
func (s *SiteService) ListByCategory(categoryID, userID uuid.UUID) ([]models.Site, error) {
	sites, err := s.siteRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	// Filter by user ownership
	var userSites []models.Site
	for _, site := range sites {
		if site.UserID == userID {
			userSites = append(userSites, site)
		}
	}

	return userSites, nil
}

// UpdateSiteInput represents site update input
type UpdateSiteInput struct {
	Name        string `json:"name" binding:"required,max=200"`
	URL         string `json:"url" binding:"required,url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

// Update updates a site
func (s *SiteService) Update(id, userID uuid.UUID, input UpdateSiteInput) (*models.Site, error) {
	// Get site and check ownership
	site, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	site.Name = input.Name
	site.URL = input.URL
	site.Description = input.Description
	site.Icon = input.Icon
	site.SortOrder = input.SortOrder

	if err := s.siteRepo.Update(site); err != nil {
		return nil, err
	}

	return site, nil
}

// Delete deletes a site
func (s *SiteService) Delete(id, userID uuid.UUID) error {
	// Check ownership first
	_, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	return s.siteRepo.Delete(id)
}

// Move moves a site to another category
func (s *SiteService) Move(id, newCategoryID, userID uuid.UUID) error {
	// Check ownership first
	_, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	return s.siteRepo.MoveSite(id, newCategoryID)
}

// Search searches for sites
func (s *SiteService) Search(userID uuid.UUID, query string) ([]models.Site, error) {
	return s.siteRepo.Search(userID, query)
}

// BatchCreateInput represents batch site creation input
type BatchCreateInput struct {
	CategoryID uuid.UUID            `json:"category_id" binding:"required"`
	Sites      []CreateSiteInput `json:"sites" binding:"required"`
}

// BatchCreate creates multiple sites
func (s *SiteService) BatchCreate(userID uuid.UUID, input BatchCreateInput) ([]models.Site, error) {
	var createdSites []models.Site

	for _, siteInput := range input.Sites {
		siteInput.CategoryID = input.CategoryID
		site, err := s.Create(userID, siteInput)
		if err != nil {
			return nil, err
		}
		createdSites = append(createdSites, *site)
	}

	return createdSites, nil
}

// ReorderSitesInput represents batch sort order update input
type ReorderSitesInput struct {
	SiteIDs []string `json:"site_ids" binding:"required"` // Array of site IDs in new order
}

// ReorderSites updates the sort order of multiple sites
func (s *SiteService) ReorderSites(userID uuid.UUID, input ReorderSitesInput) error {
	for i, siteIDStr := range input.SiteIDs {
		siteID, err := uuid.Parse(siteIDStr)
		if err != nil {
			return fmt.Errorf("invalid site ID: %s", siteIDStr)
		}

		// Check ownership
		site, err := s.GetByID(siteID, userID)
		if err != nil {
			return fmt.Errorf("site not found or access denied: %s", siteIDStr)
		}

		// Update sort order
		site.SortOrder = i + 1
		if err := s.siteRepo.Update(site); err != nil {
			return err
		}
	}

	return nil
}

// BatchDeleteInput represents batch delete input
type BatchDeleteInput struct {
	SiteIDs []string `json:"site_ids" binding:"required"`
}

// BatchDelete deletes multiple sites
func (s *SiteService) BatchDelete(userID uuid.UUID, input BatchDeleteInput) error {
	for _, siteIDStr := range input.SiteIDs {
		siteID, err := uuid.Parse(siteIDStr)
		if err != nil {
			return fmt.Errorf("invalid site ID: %s", siteIDStr)
		}

		// Check ownership
		_, err = s.GetByID(siteID, userID)
		if err != nil {
			return fmt.Errorf("site not found or access denied: %s", siteIDStr)
		}

		// Delete site
		if err := s.siteRepo.Delete(siteID); err != nil {
			return err
		}
	}

	return nil
}

// BatchUpdateInput represents batch update input
type BatchUpdateInput struct {
	Updates []SiteUpdateItem `json:"updates" binding:"required"`
}

// SiteUpdateItem represents a single site update
type SiteUpdateItem struct {
	SiteID      string `json:"site_id" binding:"required"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// BatchUpdate updates multiple sites
func (s *SiteService) BatchUpdate(userID uuid.UUID, input BatchUpdateInput) error {
	for _, update := range input.Updates {
		siteID, err := uuid.Parse(update.SiteID)
		if err != nil {
			return fmt.Errorf("invalid site ID: %s", update.SiteID)
		}

		// Get site and check ownership
		site, err := s.GetByID(siteID, userID)
		if err != nil {
			return fmt.Errorf("site not found or access denied: %s", update.SiteID)
		}

		// Update fields if provided
		if update.Name != "" {
			site.Name = update.Name
		}
		if update.URL != "" {
			site.URL = update.URL
		}
		if update.Description != "" {
			site.Description = update.Description
		}
		if update.Icon != "" {
			site.Icon = update.Icon
		}

		// Save changes
		if err := s.siteRepo.Update(site); err != nil {
			return err
		}
	}

	return nil
}
