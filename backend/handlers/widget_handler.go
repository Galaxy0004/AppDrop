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

// WidgetHandler orchestrates HTTP request processing for Widget-related resources.
type WidgetHandler struct {
	widgetRepo *repository.WidgetRepository
	pageRepo   *repository.PageRepository
}

// NewWidgetHandler initializes and returns a new instance of WidgetHandler with its required dependencies.
func NewWidgetHandler(widgetRepo *repository.WidgetRepository, pageRepo *repository.PageRepository) *WidgetHandler {
	return &WidgetHandler{
		widgetRepo: widgetRepo,
		pageRepo:   pageRepo,
	}
}

// CreateWidget processes requests to instantiate and persist a new widget within a specific page context.
func (h *WidgetHandler) CreateWidget(c *gin.Context) {
	pageIDStr := c.Param("id")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

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

	if !models.IsValidWidgetType(req.Type) {
		c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid widget type. Must be one of: banner, product_grid, text, image, spacer"))
		return
	}

	if len(req.Config) > 0 {
		var configTest interface{}
		if err := json.Unmarshal(req.Config, &configTest); err != nil {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid JSON format for widget config"))
			return
		}
	}

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

// UpdateWidget processes requests to modify the attributes or configuration of an existing widget.
func (h *WidgetHandler) UpdateWidget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid widget ID format"))
		return
	}

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

	if req.Type != nil {
		if !models.IsValidWidgetType(*req.Type) {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Invalid widget type. Must be one of: banner, product_grid, text, image, spacer"))
			return
		}
		updates["type"] = *req.Type
	}

	if req.Position != nil {
		updates["position"] = *req.Position
	}

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

// DeleteWidget processes requests to remove a widget from its associated page.
func (h *WidgetHandler) DeleteWidget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid widget ID format"))
		return
	}

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

// ReorderWidgets processes requests to synchronously update the sequential arrangement of widgets on a page.
func (h *WidgetHandler) ReorderWidgets(c *gin.Context) {
	pageIDStr := c.Param("id")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

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

	seen := make(map[uuid.UUID]bool)
	for _, id := range req.WidgetIDs {
		if seen[id] {
			c.JSON(http.StatusBadRequest, models.NewValidationError("Duplicate widget ID in the list"))
			return
		}
		seen[id] = true
	}

	widgetCount, err := h.widgetRepo.GetWidgetCountByPageID(pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to verify widgets"))
		return
	}

	if len(req.WidgetIDs) != widgetCount {
		c.JSON(http.StatusBadRequest, models.NewValidationError("The number of widget IDs must match the total widgets on the page"))
		return
	}

	if err := h.widgetRepo.Reorder(pageID, req.WidgetIDs); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to reorder widgets: "+err.Error()))
		return
	}

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

// GetWidgets processes requests to retrieve all widgets for a page, with optional type-based filtering.
func (h *WidgetHandler) GetWidgets(c *gin.Context) {
	pageIDStr := c.Param("id")
	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Invalid page ID format"))
		return
	}

	page, err := h.pageRepo.GetByID(pageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Failed to fetch page"))
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, models.NewNotFoundError("Page not found"))
		return
	}

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
