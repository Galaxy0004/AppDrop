// Package repository provides the data access layer for the application.
// It contains implementations for interacting with the persistent data store.
package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"appdrop/models"

	"github.com/google/uuid"
)

// WidgetRepository manages database operations specifically for Widget entities.
// It provides methods for CRUD operations, reordering, and validation.
type WidgetRepository struct {
	db *sql.DB
}

// NewWidgetRepository initializes and returns a new instance of WidgetRepository.
func NewWidgetRepository(db *sql.DB) *WidgetRepository {
	return &WidgetRepository{db: db}
}

// Create persists a new Widget entity in the data store.
func (r *WidgetRepository) Create(widget *models.Widget) error {
	config := widget.Config
	if config == nil {
		config = json.RawMessage("{}")
	}

	query := `
		INSERT INTO widgets (page_id, type, position, config)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, widget.PageID, widget.Type, widget.Position, config).
		Scan(&widget.ID, &widget.CreatedAt, &widget.UpdatedAt)
}

// GetByID retrieves a single Widget entity by its unique identifier.
func (r *WidgetRepository) GetByID(id uuid.UUID) (*models.Widget, error) {
	query := `
		SELECT id, page_id, type, position, config, created_at, updated_at
		FROM widgets
		WHERE id = $1
	`
	widget := &models.Widget{}
	var configBytes []byte
	err := r.db.QueryRow(query, id).Scan(
		&widget.ID, &widget.PageID, &widget.Type, &widget.Position,
		&configBytes, &widget.CreatedAt, &widget.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	widget.Config = json.RawMessage(configBytes)
	return widget, nil
}

// GetByPageID retrieves a collection of Widget entities associated with a specific Page.
// An optional widgetType filter can be applied to narrow the results.
func (r *WidgetRepository) GetByPageID(pageID uuid.UUID, widgetType *string) ([]models.Widget, error) {
	var query string
	var args []interface{}

	if widgetType != nil && *widgetType != "" {
		query = `
			SELECT id, page_id, type, position, config, created_at, updated_at
			FROM widgets
			WHERE page_id = $1 AND type = $2
			ORDER BY position ASC
		`
		args = []interface{}{pageID, *widgetType}
	} else {
		query = `
			SELECT id, page_id, type, position, config, created_at, updated_at
			FROM widgets
			WHERE page_id = $1
			ORDER BY position ASC
		`
		args = []interface{}{pageID}
	}

	rows, err := r.db.Query(query, args...)
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

	return widgets, rows.Err()
}

// GetMaxPosition determines the highest position index currently assigned to widgets on a page.
func (r *WidgetRepository) GetMaxPosition(pageID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(MAX(position), 0) FROM widgets WHERE page_id = $1`
	var maxPosition int
	err := r.db.QueryRow(query, pageID).Scan(&maxPosition)
	return maxPosition, err
}

// Update modifies an existing Widget entity with the provided field updates.
func (r *WidgetRepository) Update(id uuid.UUID, updates map[string]interface{}) (*models.Widget, error) {
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
		UPDATE widgets
		SET %s
		WHERE id = $%d
		RETURNING id, page_id, type, position, config, created_at, updated_at
	`, setClauses, argIndex)

	widget := &models.Widget{}
	var configBytes []byte
	err := r.db.QueryRow(query, args...).Scan(
		&widget.ID, &widget.PageID, &widget.Type, &widget.Position,
		&configBytes, &widget.CreatedAt, &widget.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	widget.Config = json.RawMessage(configBytes)
	return widget, nil
}

// Delete removes a Widget entity from the data store by its ID.
func (r *WidgetRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM widgets WHERE id = $1`
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

// Reorder applies a new sequential order to a list of widgets within a specific page.
func (r *WidgetRepository) Reorder(pageID uuid.UUID, widgetIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, widgetID := range widgetIDs {
		query := `UPDATE widgets SET position = $1 WHERE id = $2 AND page_id = $3`
		result, err := tx.Exec(query, i+1, widgetID, pageID)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("widget %s not found on page %s", widgetID, pageID)
		}
	}

	return tx.Commit()
}

// GetWidgetCountByPageID returns the total number of widgets associated with a specific page.
func (r *WidgetRepository) GetWidgetCountByPageID(pageID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM widgets WHERE page_id = $1`
	var count int
	err := r.db.QueryRow(query, pageID).Scan(&count)
	return count, err
}

// CheckWidgetsBelongToPage validates that a set of widget IDs all belong to the specified page identifier.
func (r *WidgetRepository) CheckWidgetsBelongToPage(pageID uuid.UUID, widgetIDs []uuid.UUID) (bool, error) {
	if len(widgetIDs) == 0 {
		return true, nil
	}

	query := `
		SELECT COUNT(*) FROM widgets 
		WHERE page_id = $1 AND id = ANY($2)
	`

	var count int
	err := r.db.QueryRow(query, pageID, widgetIDs).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == len(widgetIDs), nil
}
