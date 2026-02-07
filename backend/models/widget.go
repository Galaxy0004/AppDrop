package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WidgetType represents the type of widget
type WidgetType string

const (
	WidgetTypeBanner      WidgetType = "banner"
	WidgetTypeProductGrid WidgetType = "product_grid"
	WidgetTypeText        WidgetType = "text"
	WidgetTypeImage       WidgetType = "image"
	WidgetTypeSpacer      WidgetType = "spacer"
)

// ValidWidgetTypes contains all valid widget types
var ValidWidgetTypes = []WidgetType{
	WidgetTypeBanner,
	WidgetTypeProductGrid,
	WidgetTypeText,
	WidgetTypeImage,
	WidgetTypeSpacer,
}

// IsValidWidgetType checks if a widget type is valid
func IsValidWidgetType(t string) bool {
	for _, validType := range ValidWidgetTypes {
		if string(validType) == t {
			return true
		}
	}
	return false
}

// Widget represents a UI component on a page
type Widget struct {
	ID        uuid.UUID       `json:"id"`
	PageID    uuid.UUID       `json:"page_id"`
	Type      string          `json:"type"`
	Position  int             `json:"position"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// CreateWidgetRequest represents the request body for creating a widget
type CreateWidgetRequest struct {
	Type     string          `json:"type" binding:"required"`
	Position int             `json:"position"`
	Config   json.RawMessage `json:"config,omitempty"`
}

// UpdateWidgetRequest represents the request body for updating a widget
type UpdateWidgetRequest struct {
	Type     *string          `json:"type,omitempty"`
	Position *int             `json:"position,omitempty"`
	Config   *json.RawMessage `json:"config,omitempty"`
}

// ReorderWidgetsRequest represents the request body for reordering widgets
type ReorderWidgetsRequest struct {
	WidgetIDs []uuid.UUID `json:"widget_ids" binding:"required"`
}

// WidgetResponse represents the API response for a widget
type WidgetResponse struct {
	ID        uuid.UUID       `json:"id"`
	PageID    uuid.UUID       `json:"page_id"`
	Type      string          `json:"type"`
	Position  int             `json:"position"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
