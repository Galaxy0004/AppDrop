package models

import (
	"time"

	"github.com/google/uuid"
)

// Page represents a high-level screen or container within the mobile application.
// It serves as a parent entity for multiple UI components (widgets).
type Page struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Route     string    `json:"route"`
	IsHome    bool      `json:"is_home"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Widgets   []Widget  `json:"widgets,omitempty"`
}

// CreatePageRequest defines the expected payload for the page creation endpoint.
type CreatePageRequest struct {
	Name   string `json:"name" binding:"required"`
	Route  string `json:"route" binding:"required"`
	IsHome bool   `json:"is_home"`
}

// UpdatePageRequest defines the expected payload for the page update endpoint, where fields are optional.
type UpdatePageRequest struct {
	Name   *string `json:"name,omitempty"`
	Route  *string `json:"route,omitempty"`
	IsHome *bool   `json:"is_home,omitempty"`
}

// PageResponse encapsulates the data returned to the client for single-page queries.
type PageResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Route     string    `json:"route"`
	IsHome    bool      `json:"is_home"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Widgets   []Widget  `json:"widgets,omitempty"`
}

// PageListResponse provides a structured wrapper for paginated collections of Page entities.
type PageListResponse struct {
	Pages      []Page `json:"pages"`
	Total      int    `json:"total"`
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"per_page,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
}
