package models

import (
	"time"

	"github.com/google/uuid"
)

// Page represents a screen in the mobile app
type Page struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Route     string     `json:"route"`
	IsHome    bool       `json:"is_home"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Widgets   []Widget   `json:"widgets,omitempty"`
}

// CreatePageRequest represents the request body for creating a page
type CreatePageRequest struct {
	Name   string `json:"name" binding:"required"`
	Route  string `json:"route" binding:"required"`
	IsHome bool   `json:"is_home"`
}

// UpdatePageRequest represents the request body for updating a page
type UpdatePageRequest struct {
	Name   *string `json:"name,omitempty"`
	Route  *string `json:"route,omitempty"`
	IsHome *bool   `json:"is_home,omitempty"`
}

// PageResponse represents the API response for a page
type PageResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Route     string    `json:"route"`
	IsHome    bool      `json:"is_home"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Widgets   []Widget  `json:"widgets,omitempty"`
}

// PageListResponse represents the API response for listing pages
type PageListResponse struct {
	Pages      []Page `json:"pages"`
	Total      int    `json:"total"`
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"per_page,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
}
