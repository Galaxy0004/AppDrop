package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"appdrop/models"
	"appdrop/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WidgetHandler handles HTTP requests for widgets
type WidgetHandler struct {
	widgetRepo *repository.WidgetRepository
	pageRepo   *repository.PageRepository
}

// NewWidgetHandler creates a new WidgetHandler
func NewWidgetHandler(widgetRepo *repository.WidgetRepository, pageRepo *repository.PageRepository) *WidgetHandler {
	return &WidgetHandler{
		widgetRepo: widgetRepo,
		pageRepo:   pageRepo,
	}
}

// CreateWidget handles POST /pages/:id/widgets
func (h *WidgetHandler) CreateWidget(c *gin.Context) {
	pageIDStr := c.Param("id")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	// Check if page exists
	page, err := h.pageRepo.GetByID(pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

	var req models.CreateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid request body: "+err.Error()))
		return
	}

	// Validate widget type
	if !models.IsValidWidgetType(req.Type) {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid widget type. Must be one of: banner, product_grid, text, image, spacer"))
		return
	}

	// Validate config JSON if provided
	if req.Config != nil && len(req.Config) > 0 {
		var configTest interface{}
		if err := json.Unmarshal(req.Config, &configTest); err != nil {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid JSON format for widget config"))
			return
		}
	}

	// Get max position if position not specified
	position := req.Position
	if position == 0 {
		maxPos, err := h.widgetRepo.GetMaxPosition(pageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to determine widget position"))
			return
		}
		position = maxPos + 1
	}

	config := req.Config
	if config == nil {
		config = json.RawMessage("{}")
	}

	widget := &models.Widget{
		PageID:   pageID,
		Type:     req.Type,
		Position: position,
		Config:   config,
	}

	if err := h.widgetRepo.Create(widget); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to create widget"))
		return
	}

	c.JSON(http.StatusCreated, widget)
}

// UpdateWidget handles PUT /widgets/:id
func (h *WidgetHandler) UpdateWidget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid widget ID format"))
		return
	}

	// Check if widget exists
	existingWidget, err := h.widgetRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch widget"))
		return
	}
	if existingWidget == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Widget not found"))
		return
	}

	var req models.UpdateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid request body: "+err.Error()))
		return
	}

	updates := make(map[string]interface{})

	// Validate and add type update
	if req.Type != nil {
		if !models.IsValidWidgetType(*req.Type) {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid widget type. Must be one of: banner, product_grid, text, image, spacer"))
			return
		}
		updates["type"] = *req.Type
	}

	// Add position update
	if req.Position != nil {
		updates["position"] = *req.Position
	}

	// Validate and add config update
	if req.Config != nil {
		var configTest interface{}
		if err := json.Unmarshal(*req.Config, &configTest); err != nil {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid JSON format for widget config"))
			return
		}
		updates["config"] = *req.Config
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, existingWidget)
		return
	}

	widget, err := h.widgetRepo.Update(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to update widget"))
		return
	}

	c.JSON(http.StatusOK, widget)
}

// DeleteWidget handles DELETE /widgets/:id
func (h *WidgetHandler) DeleteWidget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid widget ID format"))
		return
	}

	// Check if widget exists
	widget, err := h.widgetRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch widget"))
		return
	}
	if widget == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Widget not found"))
		return
	}

	if err := h.widgetRepo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.NewNotFoundError("Widget not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to delete widget"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Widget deleted successfully"})
}

// ReorderWidgets handles POST /pages/:id/widgets/reorder
func (h *WidgetHandler) ReorderWidgets(c *gin.Context) {
	pageIDStr := c.Param("id")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	// Check if page exists
	page, err := h.pageRepo.GetByID(pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

	var req models.ReorderWidgetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid request body: "+err.Error()))
		return
	}

	if len(req.WidgetIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Widget IDs array cannot be empty"))
		return
	}

	// Check for duplicates
	seen := make(map[uuid.UUID]bool)
	for _, id := range req.WidgetIDs {
		if seen[id] {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Duplicate widget ID in the list"))
			return
		}
		seen[id] = true
	}

	// Verify all widgets belong to this page
	widgetCount, err := h.widgetRepo.GetWidgetCountByPageID(pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to verify widgets"))
		return
	}

	if len(req.WidgetIDs) != widgetCount {
		c.JSON(http.StatusBadRequest, models.NewValidationError("The number of widget IDs must match the total widgets on the page"))
		return
	}

	// Reorder widgets
	if err := h.widgetRepo.Reorder(pageID, req.WidgetIDs); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to reorder widgets: "+err.Error()))
		return
	}

	// Get updated widgets
	widgets, err := h.widgetRepo.GetByPageID(pageID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch updated widgets"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Widgets reordered successfully",
		"widgets": widgets,
	})
}

// GetWidgets handles GET /pages/:id/widgets (bonus: with type filter)
func (h *WidgetHandler) GetWidgets(c *gin.Context) {
	pageIDStr := c.Param("id")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	// Check if page exists
	page, err := h.pageRepo.GetByID(pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

	// Get widget type filter
	var widgetType *string
	if t := c.Query("type"); t != "" {
		if !models.IsValidWidgetType(t) {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid widget type filter"))
			return
		}
		widgetType = &t
	}

	widgets, err := h.widgetRepo.GetByPageID(pageID, widgetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch widgets"))
		return
	}

	if widgets == nil {
		widgets = []models.Widget{}
	}

	c.JSON(http.StatusOK, gin.H{
		"widgets": widgets,
		"total":   len(widgets),
	})
}
