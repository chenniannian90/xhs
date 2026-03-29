package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/navhub/internal/services"
)

// SiteHandler handles site HTTP requests
type SiteHandler struct {
	siteService *services.SiteService
}

// NewSiteHandler creates a new site handler
func NewSiteHandler(siteService *services.SiteService) *SiteHandler {
	return &SiteHandler{
		siteService: siteService,
	}
}

// ListSites gets all sites for the authenticated user
func (h *SiteHandler) ListSites(c *gin.Context) {
	userID, _ := c.Get("userID")

	sites, err := h.siteService.ListByUser(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Sites retrieved",
		Data:    sites,
	})
}

// CreateSite creates a new site
func (h *SiteHandler) CreateSite(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.CreateSiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	site, err := h.siteService.Create(userID.(uuid.UUID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    201,
		Message: "Site created",
		Data:    site,
	})
}

// GetSite gets a site by ID
func (h *SiteHandler) GetSite(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid site ID"})
		return
	}

	site, err := h.siteService.GetByID(id, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Site not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Site retrieved",
		Data:    site,
	})
}

// UpdateSite updates a site
func (h *SiteHandler) UpdateSite(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid site ID"})
		return
	}

	var input services.UpdateSiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	site, err := h.siteService.Update(id, userID.(uuid.UUID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Site updated",
		Data:    site,
	})
}

// DeleteSite deletes a site
func (h *SiteHandler) DeleteSite(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid site ID"})
		return
	}

	if err := h.siteService.Delete(id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Site deleted",
	})
}

// MoveSite moves a site to another category
func (h *SiteHandler) MoveSite(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid site ID"})
		return
	}

	var req struct {
		CategoryID uuid.UUID `json:"category_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.siteService.Move(id, req.CategoryID, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Site moved",
	})
}

// SearchSites searches for sites
func (h *SiteHandler) SearchSites(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Query parameter 'q' is required"})
		return
	}

	sites, err := h.siteService.Search(userID.(uuid.UUID), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Sites retrieved",
		Data:    sites,
	})
}

// BatchCreateSites creates multiple sites at once
func (h *SiteHandler) BatchCreateSites(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input services.BatchCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	sites, err := h.siteService.BatchCreate(userID.(uuid.UUID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    201,
		Message: "Sites created",
		Data:    sites,
	})
}

// GlobalSearchSites performs a global site search (for search page)
func (h *SiteHandler) GlobalSearchSites(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")

	if query == "" {
		// Return empty results if no query
		c.JSON(http.StatusOK, SuccessResponse{
			Code:    200,
			Message: "Sites retrieved",
			Data:    []interface{}{},
		})
		return
	}

	sites, err := h.siteService.Search(userID.(uuid.UUID), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Search results",
		Data:    sites,
	})
}

// ReorderSites handles batch site reorder
func (h *SiteHandler) ReorderSites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	var input services.ReorderSitesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(input.SiteIDs) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No sites to reorder"})
		return
	}

	if err := h.siteService.ReorderSites(userID.(uuid.UUID), input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Sites reordered successfully",
	})
}

// BatchDelete handles batch site deletion
func (h *SiteHandler) BatchDelete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	var input services.BatchDeleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(input.SiteIDs) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No sites to delete"})
		return
	}

	if err := h.siteService.BatchDelete(userID.(uuid.UUID), input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Sites deleted successfully",
	})
}

// BatchUpdate handles batch site update
func (h *SiteHandler) BatchUpdate(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Not authenticated"})
		return
	}

	var input services.BatchUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(input.Updates) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No updates provided"})
		return
	}

	if err := h.siteService.BatchUpdate(userID.(uuid.UUID), input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "Sites updated successfully",
	})
}
