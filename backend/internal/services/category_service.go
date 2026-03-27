package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/models"
	"github.com/yourusername/navhub/internal/repositories"
	"github.com/yourusername/navhub/pkg/validator"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
	siteRepo     *repositories.SiteRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo *repositories.CategoryRepository, siteRepo *repositories.SiteRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		siteRepo:     siteRepo,
	}
}

// CreateCategoryInput represents category creation input
type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

// Create creates a new category
func (s *CategoryService) Create(userID uuid.UUID, input CreateCategoryInput) (*models.Category, error) {
	// Validate input
	v := validator.NewValidator()
	v.ValidateName("name", input.Name, 1, 100)
	v.ValidateDescription("description", input.Description, 500)
	v.ValidateIcon("icon", input.Icon)
	v.ValidateSortOrder("sort_order", input.SortOrder)

	if v.HasErrors() {
		return nil, v.Error()
	}

	// Check if category name already exists for user
	categories, err := s.categoryRepo.FindByUserID(userID)
	if err == nil {
		existingNames := make([]string, len(categories))
		for i, cat := range categories {
			existingNames[i] = cat.Name
		}
		if !validator.IsUnique(input.Name, existingNames) {
			return nil, errors.New("category with this name already exists")
		}
	}

	category := &models.Category{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Icon:        input.Icon,
		SortOrder:   input.SortOrder,
		IsPublic:    false,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetByID gets a category by ID
func (s *CategoryService) GetByID(id, userID uuid.UUID) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if category.UserID != userID {
		return nil, errors.New("access denied")
	}

	return category, nil
}

// ListByUser gets all categories for a user
func (s *CategoryService) ListByUser(userID uuid.UUID) ([]models.Category, error) {
	return s.categoryRepo.FindByUserID(userID)
}

// UpdateCategoryInput represents category update input
type UpdateCategoryInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

// Update updates a category
func (s *CategoryService) Update(id, userID uuid.UUID, input UpdateCategoryInput) (*models.Category, error) {
	// Get category and check ownership
	category, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	category.Name = input.Name
	category.Description = input.Description
	category.Icon = input.Icon
	category.SortOrder = input.SortOrder

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete deletes a category
func (s *CategoryService) Delete(id, userID uuid.UUID) error {
	// Check ownership first
	_, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(id)
}

// Share generates a share token for a category
func (s *CategoryService) Share(id, userID uuid.UUID) (string, error) {
	category, err := s.GetByID(id, userID)
	if err != nil {
		return "", err
	}

	// Generate share token
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	category.IsPublic = true
	category.ShareToken = token

	if err := s.categoryRepo.Update(category); err != nil {
		return "", err
	}

	return token, nil
}

// Unshare removes public sharing from a category
func (s *CategoryService) Unshare(id, userID uuid.UUID) error {
	category, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}

	category.IsPublic = false
	category.ShareToken = ""

	return s.categoryRepo.Update(category)
}

// GetByShareToken gets a public category by share token
func (s *CategoryService) GetByShareToken(token string) (*models.Category, error) {
	return s.categoryRepo.FindByShareToken(token)
}

// ExportToCSV exports all user data to CSV format
func (s *CategoryService) ExportToCSV(userID uuid.UUID) ([]byte, error) {
	// Get all categories for user
	categories, err := s.ListByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	fmt.Printf("📊 [DEBUG] Found %d categories for user %s\n", len(categories), userID)

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	defer writer.Flush()

	// Write CSV header
	header := []string{"type", "id", "name", "description", "icon", "sort_order", "site_id", "site_name", "site_url", "site_description", "site_icon", "site_sort_order"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	fmt.Printf("📊 [DEBUG] CSV header written, buffer size: %d\n", buffer.Len())

	// Write category data
	for _, category := range categories {
		// Write category row
		categoryRow := []string{
			"category",
			category.ID.String(),
			category.Name,
			category.Description,
			category.Icon,
			strconv.Itoa(category.SortOrder),
			"",
			"",
			"",
			"",
			"",
			"",
		}
		if err := writer.Write(categoryRow); err != nil {
			return nil, fmt.Errorf("failed to write category row: %w", err)
		}

		// Write site data for this category
		for _, site := range category.Sites {
			siteRow := []string{
				"site",
				category.ID.String(),
				"",
				"",
				"",
				"",
				"",
				site.ID.String(),
				site.Name,
				site.URL,
				site.Description,
				site.Icon,
				strconv.Itoa(site.SortOrder),
			}
			if err := writer.Write(siteRow); err != nil {
				return nil, fmt.Errorf("failed to write site row: %w", err)
			}
		}
	}

	fmt.Printf("📊 [DEBUG] Final CSV buffer size: %d bytes\n", buffer.Len())
	return buffer.Bytes(), nil
}

// ImportInput represents CSV import input
type ImportInput struct {
	CategoryID uuid.UUID `json:"category_id"` // Optional: import into specific category
}

// ImportResult represents the result of CSV import
type ImportResult struct {
	CategoriesCreated int       `json:"categories_created"`
	SitesCreated      int          `json:"sites_created"`
	Errors            []string      `json:"errors"`
}

// ImportFromCSV imports data from CSV format
func (s *CategoryService) ImportFromCSV(userID uuid.UUID, csvData []byte, input ImportInput) (*ImportResult, error) {
	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return &ImportResult{}, nil
	}

	result := &ImportResult{}
	categories := make(map[string]*models.Category)

	// Process CSV records
	for i, record := range records {
		if len(record) < 10 {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Invalid column count", i+1))
			continue
		}

		recordType := record[0]

		if recordType == "category" {
			// Import category
			category := &models.Category{
				UserID:      userID,
				Name:        record[2],
				Description: record[3],
				Icon:        record[4],
				SortOrder:   parseInt(record[5]),
				IsPublic:    false,
			}

			if err := s.categoryRepo.Create(category); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to create category: %w", i+1, err))
			} else {
				result.CategoriesCreated++
				categories[record[1]] = category
			}

		} else if recordType == "site" {
			// Import site
			categoryID := input.CategoryID
			if categoryID == uuid.Nil {
				// Use category from CSV
				if cat, exists := categories[record[1]]; exists {
					categoryID = cat.ID
				} else {
					result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Category not found", i+1))
					continue
				}
			}

			site := &models.Site{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        record[7],
				URL:         record[8],
				Description: record[9],
				Icon:        record[10],
				SortOrder:   parseInt(record[11]),
			}

			if err := s.siteRepo.Create(site); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to create site: %w", i+1, err))
			} else {
				result.SitesCreated++
			}
		}
	}

	return result, nil
}

// parseInt safely parses a string to int
func parseInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

// ReorderCategoriesInput represents batch sort order update input
type ReorderCategoriesInput struct {
	CategoryIDs []string `json:"category_ids" binding:"required"` // Array of category IDs in new order
}

// ReorderCategories updates the sort order of multiple categories
func (s *CategoryService) ReorderCategories(userID uuid.UUID, input ReorderCategoriesInput) error {
	for i, categoryIDStr := range input.CategoryIDs {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return fmt.Errorf("invalid category ID: %s", categoryIDStr)
		}

		// Check ownership
		category, err := s.GetByID(categoryID, userID)
		if err != nil {
			return fmt.Errorf("category not found or access denied: %s", categoryIDStr)
		}

		// Update sort order
		category.SortOrder = i + 1
		if err := s.categoryRepo.Update(category); err != nil {
			return err
		}
	}

	return nil
}

// BatchDeleteCategoriesInput represents batch category delete input
type BatchDeleteCategoriesInput struct {
	CategoryIDs []string `json:"category_ids" binding:"required"`
}

// BatchDeleteCategories deletes multiple categories
func (s *CategoryService) BatchDeleteCategories(userID uuid.UUID, input BatchDeleteCategoriesInput) error {
	for _, categoryIDStr := range input.CategoryIDs {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return fmt.Errorf("invalid category ID: %s", categoryIDStr)
		}

		// Check ownership
		_, err = s.GetByID(categoryID, userID)
		if err != nil {
			return fmt.Errorf("category not found or access denied: %s", categoryIDStr)
		}

		// Delete category
		if err := s.categoryRepo.Delete(categoryID); err != nil {
			return err
		}
	}

	return nil
}

// BatchUpdateCategoriesInput represents batch category update input
type BatchUpdateCategoriesInput struct {
	Updates []CategoryUpdateItem `json:"updates" binding:"required"`
}

// CategoryUpdateItem represents a single category update
type CategoryUpdateItem struct {
	CategoryID  string `json:"category_id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// BatchUpdateCategories updates multiple categories
func (s *CategoryService) BatchUpdateCategories(userID uuid.UUID, input BatchUpdateCategoriesInput) error {
	for _, update := range input.Updates {
		categoryID, err := uuid.Parse(update.CategoryID)
		if err != nil {
			return fmt.Errorf("invalid category ID: %s", update.CategoryID)
		}

		// Get category and check ownership
		category, err := s.GetByID(categoryID, userID)
		if err != nil {
			return fmt.Errorf("category not found or access denied: %s", update.CategoryID)
		}

		// Update fields if provided
		if update.Name != "" {
			category.Name = update.Name
		}
		if update.Description != "" {
			category.Description = update.Description
		}
		if update.Icon != "" {
			category.Icon = update.Icon
		}

		// Save changes
		if err := s.categoryRepo.Update(category); err != nil {
			return err
		}
	}

	return nil
}

// ExportSingleCategoryToCSV exports a single category with its sites to CSV
func (s *CategoryService) ExportSingleCategoryToCSV(category *models.Category) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	defer writer.Flush()

	// Write CSV header
	header := []string{"type", "category_id", "category_name", "site_id", "site_name", "site_url", "site_description", "site_icon", "site_sort_order"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write category row
	categoryRow := []string{
		"category",
		category.ID.String(),
		category.Name,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}
	if err := writer.Write(categoryRow); err != nil {
		return nil, fmt.Errorf("failed to write category row: %w", err)
	}

	// Write site data for this category
	for _, site := range category.Sites {
		siteRow := []string{
			"site",
			category.ID.String(),
			category.Name,
			site.ID.String(),
			site.Name,
			site.URL,
			site.Description,
			site.Icon,
			strconv.Itoa(site.SortOrder),
		}
		if err := writer.Write(siteRow); err != nil {
			return nil, fmt.Errorf("failed to write site row: %w", err)
		}
	}

	return buffer.Bytes(), nil
}
