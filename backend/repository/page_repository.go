package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"appdrop/models"

	"github.com/google/uuid"
)

// PageRepository handles database operations for pages
type PageRepository struct {
	db *sql.DB
}

// NewPageRepository creates a new PageRepository
func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{db: db}
}

// Create creates a new page in the database
func (r *PageRepository) Create(page *models.Page) error {
	query := `
		INSERT INTO pages (name, route, is_home)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, page.Name, page.Route, page.IsHome).
		Scan(&page.ID, &page.CreatedAt, &page.UpdatedAt)
}

// GetByID retrieves a page by its ID
func (r *PageRepository) GetByID(id uuid.UUID) (*models.Page, error) {
	query := `
		SELECT id, name, route, is_home, created_at, updated_at
		FROM pages
		WHERE id = $1
	`
	page := &models.Page{}
	err := r.db.QueryRow(query, id).Scan(
		&page.ID, &page.Name, &page.Route, &page.IsHome,
		&page.CreatedAt, &page.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return page, nil
}

// GetByRoute retrieves a page by its route
func (r *PageRepository) GetByRoute(route string) (*models.Page, error) {
	query := `
		SELECT id, name, route, is_home, created_at, updated_at
		FROM pages
		WHERE route = $1
	`
	page := &models.Page{}
	err := r.db.QueryRow(query, route).Scan(
		&page.ID, &page.Name, &page.Route, &page.IsHome,
		&page.CreatedAt, &page.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return page, nil
}

// GetHomePage retrieves the home page
func (r *PageRepository) GetHomePage() (*models.Page, error) {
	query := `
		SELECT id, name, route, is_home, created_at, updated_at
		FROM pages
		WHERE is_home = TRUE
		LIMIT 1
	`
	page := &models.Page{}
	err := r.db.QueryRow(query).Scan(
		&page.ID, &page.Name, &page.Route, &page.IsHome,
		&page.CreatedAt, &page.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return page, nil
}

// GetAll retrieves all pages with pagination
func (r *PageRepository) GetAll(page, perPage int) ([]models.Page, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM pages`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get pages with pagination
	offset := (page - 1) * perPage
	query := `
		SELECT id, name, route, is_home, created_at, updated_at
		FROM pages
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Route, &p.IsHome,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		pages = append(pages, p)
	}

	return pages, total, rows.Err()
}

// Update updates a page in the database
func (r *PageRepository) Update(id uuid.UUID, updates map[string]interface{}) (*models.Page, error) {
	// Build dynamic update query
	setClauses := ""
	args := []interface{}{}
	argIndex := 1

	for key, value := range updates {
		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += fmt.Sprintf("%s = $%d", key, argIndex)
		args = append(args, value)
		argIndex++
	}

	if setClauses == "" {
		return r.GetByID(id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE pages
		SET %s
		WHERE id = $%d
		RETURNING id, name, route, is_home, created_at, updated_at
	`, setClauses, argIndex)

	page := &models.Page{}
	err := r.db.QueryRow(query, args...).Scan(
		&page.ID, &page.Name, &page.Route, &page.IsHome,
		&page.CreatedAt, &page.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return page, nil
}

// Delete deletes a page from the database
func (r *PageRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM pages WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UnsetHomePage unsets the current home page
func (r *PageRepository) UnsetHomePage() error {
	query := `UPDATE pages SET is_home = FALSE WHERE is_home = TRUE`
	_, err := r.db.Exec(query)
	return err
}

// GetByIDWithWidgets retrieves a page with its widgets
func (r *PageRepository) GetByIDWithWidgets(id uuid.UUID) (*models.Page, error) {
	page, err := r.GetByID(id)
	if err != nil || page == nil {
		return page, err
	}

	// Get widgets for this page
	query := `
		SELECT id, page_id, type, position, config, created_at, updated_at
		FROM widgets
		WHERE page_id = $1
		ORDER BY position ASC
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var widgets []models.Widget
	for rows.Next() {
		var w models.Widget
		var configBytes []byte
		if err := rows.Scan(
			&w.ID, &w.PageID, &w.Type, &w.Position,
			&configBytes, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		w.Config = json.RawMessage(configBytes)
		widgets = append(widgets, w)
	}

	page.Widgets = widgets
	return page, rows.Err()
}

// CheckRouteExists checks if a route exists (excluding a specific page ID)
func (r *PageRepository) CheckRouteExists(route string, excludeID *uuid.UUID) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT EXISTS(SELECT 1 FROM pages WHERE route = $1 AND id != $2)`
		args = []interface{}{route, *excludeID}
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM pages WHERE route = $1)`
		args = []interface{}{route}
	}

	var exists bool
	err := r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}
