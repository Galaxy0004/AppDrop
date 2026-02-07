package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"appdrop/models"
	"appdrop/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PageHandler handles HTTP requests for pages
type PageHandler struct {
	pageRepo   *repository.PageRepository
	widgetRepo *repository.WidgetRepository
}

// NewPageHandler creates a new PageHandler
func NewPageHandler(pageRepo *repository.PageRepository, widgetRepo *repository.WidgetRepository) *PageHandler {
	return &PageHandler{
		pageRepo:   pageRepo,
		widgetRepo: widgetRepo,
	}
}

// ListPages handles GET /pages
func (h *PageHandler) ListPages(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	perPage := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	pages, total, err := h.pageRepo.GetAll(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch pages"))
		return
	}

	if pages == nil {
		pages = []models.Page{}
	}

	totalPages := (total + perPage - 1) / perPage

	c.JSON(http.StatusOK, models.PageListResponse{
		Pages:      pages,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetPage handles GET /pages/:id
func (h *PageHandler) GetPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	page, err := h.pageRepo.GetByIDWithWidgets(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

	c.JSON(http.StatusOK, page)
}

// CreatePage handles POST /pages
func (h *PageHandler) CreatePage(c *gin.Context) {
	var req models.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid request body: "+err.Error()))
		return
	}

	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Page name is required and cannot be empty"))
		return
	}

	// Validate route
	if strings.TrimSpace(req.Route) == "" {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Page route is required and cannot be empty"))
		return
	}

	// Check if route already exists
	exists, err := h.pageRepo.CheckRouteExists(req.Route, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to check route"))
		return
	}
	if exists {
		c.JSON(http.StatusConflict, models.NewConflictError("Page route already exists"))
		return
	}

	// If this is set as home page, unset the current home page
	if req.IsHome {
		if err := h.pageRepo.UnsetHomePage(); err != nil {
			c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to update home page"))
			return
		}
	}

	page := &models.Page{
		Name:   strings.TrimSpace(req.Name),
		Route:  strings.TrimSpace(req.Route),
		IsHome: req.IsHome,
	}

	if err := h.pageRepo.Create(page); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to create page"))
		return
	}

	c.JSON(http.StatusCreated, page)
}

// UpdatePage handles PUT /pages/:id
func (h *PageHandler) UpdatePage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	// Check if page exists
	existingPage, err := h.pageRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if existingPage == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

	var req models.UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid request body: "+err.Error()))
		return
	}

	updates := make(map[string]interface{})

	// Validate and add name update
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Page name cannot be empty"))
			return
		}
		updates["name"] = strings.TrimSpace(*req.Name)
	}

	// Validate and add route update
	if req.Route != nil {
		if strings.TrimSpace(*req.Route) == "" {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Page route cannot be empty"))
			return
		}
		// Check if route already exists
		exists, err := h.pageRepo.CheckRouteExists(*req.Route, &id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to check route"))
			return
		}
		if exists {
			c.JSON(http.StatusConflict, models.NewConflictError("Page route already exists"))
			return
		}
		updates["route"] = strings.TrimSpace(*req.Route)
	}

	// Handle is_home update
	if req.IsHome != nil {
		if *req.IsHome && !existingPage.IsHome {
			// Setting as home page - unset the current home page first
			if err := h.pageRepo.UnsetHomePage(); err != nil {
				c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to update home page"))
				return
			}
		}
		updates["is_home"] = *req.IsHome
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, existingPage)
		return
	}

	page, err := h.pageRepo.Update(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to update page"))
		return
	}

	c.JSON(http.StatusOK, page)
}

// DeletePage handles DELETE /pages/:id
func (h *PageHandler) DeletePage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	// Check if page exists
	page, err := h.pageRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

	// Cannot delete home page
	if page.IsHome {
		c.JSON(http.StatusConflict, models.NewConflictError("Cannot delete the home page. Set another page as home first."))
		return
	}

	// Delete the page (widgets will be cascade deleted)
	if err := h.pageRepo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to delete page"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
