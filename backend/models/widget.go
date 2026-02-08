package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WidgetType defines a custom string type for categorizing different UI component implementations.
type WidgetType string

const (
	// WidgetTypeBanner represents a full-width promotional image or slider component.
	WidgetTypeBanner WidgetType = "banner"
	// WidgetTypeProductGrid represents a layout component for displaying multiple products in a grid.
	WidgetTypeProductGrid WidgetType = "product_grid"
	// WidgetTypeText represents a simple text-based content component.
	WidgetTypeText WidgetType = "text"
	// WidgetTypeImage represents a general-purpose image component.
	WidgetTypeImage WidgetType = "image"
	// WidgetTypeSpacer represents a structural component used to add vertical or horizontal whitespace.
	WidgetTypeSpacer WidgetType = "spacer"
)

// ValidWidgetTypes maintains an authoritative list of all supported WidgetType values for validation purposes.
var ValidWidgetTypes = []WidgetType{
	WidgetTypeBanner,
	WidgetTypeProductGrid,
	WidgetTypeText,
	WidgetTypeImage,
	WidgetTypeSpacer,
}

// IsValidWidgetType performs a validation check to ensure a given string corresponds to a recognized WidgetType.
func IsValidWidgetType(t string) bool {
	for _, validType := range ValidWidgetTypes {
		if string(validType) == t {
			return true
		}
	}
	return false
}

// Widget represents an individual UI configuration element belonging to a specific Page.
type Widget struct {
	ID        uuid.UUID       `json:"id"`
	PageID    uuid.UUID       `json:"page_id"`
	Type      string          `json:"type"`
	Position  int             `json:"position"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// CreateWidgetRequest defines the data required to instantiate and persist a new Widget.
type CreateWidgetRequest struct {
	Type     string          `json:"type" binding:"required"`
	Position int             `json:"position"`
	Config   json.RawMessage `json:"config,omitempty"`
}

// UpdateWidgetRequest defines the structure for partially updating an existing Widget's configuration.
type UpdateWidgetRequest struct {
	Type     *string          `json:"type,omitempty"`
	Position *int             `json:"position,omitempty"`
	Config   *json.RawMessage `json:"config,omitempty"`
}

// ReorderWidgetsRequest defines the payload for updating the sequential ordering of widgets on a page.
type ReorderWidgetsRequest struct {
	WidgetIDs []uuid.UUID `json:"widget_ids" binding:"required"`
}

// WidgetResponse encapsulate the Widget entity data for API consumption.
type WidgetResponse struct {
	ID        uuid.UUID       `json:"id"`
	PageID    uuid.UUID       `json:"page_id"`
	Type      string          `json:"type"`
	Position  int             `json:"position"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
