package models

import (
	"testing"
)

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

func TestNewErrorResponse(t *testing.T) {
	response := NewErrorResponse(ErrorCodeValidation, "test message")
	
	if response.Error.Code != ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrorCodeValidation, response.Error.Code)
	}
	if response.Error.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", response.Error.Message)
	}
}

func TestNewValidationError(t *testing.T) {
	response := NewValidationError("validation failed")
	
	if response.Error.Code != ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrorCodeValidation, response.Error.Code)
	}
	if response.Error.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", response.Error.Message)
	}
}

func TestNewNotFoundError(t *testing.T) {
	response := NewNotFoundError("not found")
	
	if response.Error.Code != ErrorCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrorCodeNotFound, response.Error.Code)
	}
}

func TestNewConflictError(t *testing.T) {
	response := NewConflictError("conflict")
	
	if response.Error.Code != ErrorCodeConflict {
		t.Errorf("Expected code %s, got %s", ErrorCodeConflict, response.Error.Code)
	}
}

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
