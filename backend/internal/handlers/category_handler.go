package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/services"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// ListCategories gets all categories for the authenticated user
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	userID, _ := c.Get("userID")

	categories, err := h.categoryService.ListByUser(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Categories retrieved",
		Data:    categories,
	})
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.Create(userID.(uuid.UUID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    201,
		Message: "Category created",
		Data:    category,
	})
}

// GetCategory gets a category by ID
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetByID(id, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Category not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Category retrieved",
		Data:    category,
	})
}

// UpdateCategory updates a category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid category ID"})
		return
	}

	var input services.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.categoryService.Update(id, userID.(uuid.UUID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Category updated",
		Data:    category,
	})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid category ID"})
		return
	}

	if err := h.categoryService.Delete(id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Category deleted",
	})
}

// ShareCategory generates a share link for a category
func (h *CategoryHandler) ShareCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid category ID"})
		return
	}

	token, err := h.categoryService.Share(id, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	shareURL := "http://localhost:5173/shared/" + token

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Category shared",
		Data:    gin.H{"share_url": shareURL, "token": token},
	})
}

// UnshareCategory removes public sharing from a category
func (h *CategoryHandler) UnshareCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid category ID"})
		return
	}

	if err := h.categoryService.Unshare(id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Category unshared",
	})
}

// GetSharedCategory gets a public category by share token
func (h *CategoryHandler) GetSharedCategory(c *gin.Context) {
	token := c.Param("token")

	category, err := h.categoryService.GetByShareToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Shared category not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Shared category retrieved",
		Data:    category,
	})
}

// SearchCategories searches for categories
func (h *CategoryHandler) SearchCategories(c *gin.Context) {
	userID, _ := c.Get("userID")

	categories, err := h.categoryService.ListByUser(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// TODO: Implement search logic
	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Categories retrieved",
		Data:    categories,
	})
}

// ExportCategories exports all user data to CSV
func (h *CategoryHandler) ExportCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	csvData, err := h.categoryService.ExportToCSV(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if len(csvData) == 0 {
		c.JSON(http.StatusOK, SuccessResponse{
			Code:    200,
			Message: "No data to export",
			Data:    nil,
		})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=navhub-export.csv")
	c.Data(http.StatusOK, "text/csv", csvData)
}

// ExportCategory exports a single category to CSV
func (h *CategoryHandler) ExportCategory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	categoryID := c.Param("id")
	id, err := uuid.Parse(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetByID(id, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Category not found"})
		return
	}

	// Generate CSV for single category
	csvData, err := h.categoryService.ExportSingleCategoryToCSV(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", category.Name))
	c.Data(http.StatusOK, "text/csv", csvData)
}

// ImportCategories imports data from CSV
func (h *CategoryHandler) ImportCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	// Read CSV data from request body
	csvData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to read CSV data"})
		return
	}

	if len(csvData) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "CSV data is empty"})
		return
	}

	var input services.ImportInput
	if c.ShouldBindJSON(&input) != nil {
		// If no category specified, import into first available category
		// or import logic handled in service
	}

	result, err := h.categoryService.ImportFromCSV(userID.(uuid.UUID), csvData, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Import completed",
		Data:    result,
	})
}

// ReorderCategories handles batch category reorder
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	var input services.ReorderCategoriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(input.CategoryIDs) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No categories to reorder"})
		return
	}

	if err := h.categoryService.ReorderCategories(userID.(uuid.UUID), input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Categories reordered successfully",
	})
}

// BatchDeleteCategories handles batch category deletion
func (h *CategoryHandler) BatchDeleteCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	var input services.BatchDeleteCategoriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(input.CategoryIDs) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No categories to delete"})
		return
	}

	if err := h.categoryService.BatchDeleteCategories(userID.(uuid.UUID), input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Categories deleted successfully",
	})
}

// BatchUpdateCategories handles batch category update
func (h *CategoryHandler) BatchUpdateCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	var input services.BatchUpdateCategoriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(input.Updates) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No updates provided"})
		return
	}

	if err := h.categoryService.BatchUpdateCategories(userID.(uuid.UUID), input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Categories updated successfully",
	})
}
