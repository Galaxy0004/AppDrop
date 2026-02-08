// Package models contains unit tests for verifying the integrity of data models and validation logic.
package models

import (
	"testing"
)

// TestIsValidWidgetType verifies that the widget type validation logic correctly identifies
// supported and unsupported widget type strings across various edge cases.
func TestIsValidWidgetType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid banner", "banner", true},
		{"valid product_grid", "product_grid", true},
		{"valid text", "text", true},
		{"valid image", "image", true},
		{"valid spacer", "spacer", true},
		{"invalid type", "invalid", false},
		{"empty string", "", false},
		{"mixed case", "Banner", false},
		{"with spaces", " banner ", false},
		{"numbers", "123", false},
		{"special chars", "banner!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidWidgetType(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidWidgetType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNewErrorResponse ensures that generalized error responses are initialized with the
// correct code and message parameters.
func TestNewErrorResponse(t *testing.T) {
	response := NewErrorResponse(ErrorCodeValidation, "test message")

	if response.Error.Code != ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrorCodeValidation, response.Error.Code)
	}
	if response.Error.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", response.Error.Message)
	}
}

// TestNewValidationError verifies the shorthand constructor for validation errors.
func TestNewValidationError(t *testing.T) {
	response := NewValidationError("validation failed")

	if response.Error.Code != ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrorCodeValidation, response.Error.Code)
	}
	if response.Error.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", response.Error.Message)
	}
}

// TestNewNotFoundError verifies the shorthand constructor for resource not found errors.
func TestNewNotFoundError(t *testing.T) {
	response := NewNotFoundError("not found")

	if response.Error.Code != ErrorCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrorCodeNotFound, response.Error.Code)
	}
}

// TestNewConflictError verifies the shorthand constructor for resource conflict errors.
func TestNewConflictError(t *testing.T) {
	response := NewConflictError("conflict")

	if response.Error.Code != ErrorCodeConflict {
		t.Errorf("Expected code %s, got %s", ErrorCodeConflict, response.Error.Code)
	}
}

// TestValidWidgetTypes confirms that the global list of valid widget types is complete
// and ordered correctly as per the architectural design.
func TestValidWidgetTypes(t *testing.T) {
	expectedTypes := []WidgetType{
		WidgetTypeBanner,
		WidgetTypeProductGrid,
		WidgetTypeText,
		WidgetTypeImage,
		WidgetTypeSpacer,
	}

	if len(ValidWidgetTypes) != len(expectedTypes) {
		t.Errorf("Expected %d widget types, got %d", len(expectedTypes), len(ValidWidgetTypes))
	}

	for i, wt := range expectedTypes {
		if ValidWidgetTypes[i] != wt {
			t.Errorf("Expected widget type %s at index %d, got %s", wt, i, ValidWidgetTypes[i])
		}
	}
}
