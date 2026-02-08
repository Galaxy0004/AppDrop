// Package models defines the core data structures and types used across the application.
// It includes entity definitions, request/response models, and error handling structures.
package models

// ErrorCode defines a standardized set of string constants representing specific application error states.
type ErrorCode string

const (
	// ErrorCodeValidation indicates that the provided input failed business logic or schema validation.
	ErrorCodeValidation ErrorCode = "VALIDATION_ERROR"
	// ErrorCodeNotFound indicates that the requested resource could not be located in the system.
	ErrorCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrorCodeConflict indicates that the operation would result in a state conflict, such as a duplicate entry.
	ErrorCodeConflict ErrorCode = "CONFLICT"
	// ErrorCodeInternalServer indicates an unexpected failure within the server infrastructure.
	ErrorCodeInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	// ErrorCodeBadRequest indicates that the client request was malformed or syntactically incorrect.
	ErrorCodeBadRequest ErrorCode = "BAD_REQUEST"
)

// APIError encapsulates the details of a specific error, including a machine-readable code and a human-readable message.
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// ErrorResponse defines the standard structure for all error responses returned by the API.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// NewErrorResponse initializes a new instance of ErrorResponse with the specified error code and descriptive message.
func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	}
}

// NewValidationError initializes an ErrorResponse specifically for validation-related failures.
func NewValidationError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeValidation, message)
}

// NewNotFoundError initializes an ErrorResponse specifically for resource-not-found failures.
func NewNotFoundError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeNotFound, message)
}

// NewConflictError initializes an ErrorResponse specifically for data conflict failures.
func NewConflictError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeConflict, message)
}

// NewInternalServerError initializes an ErrorResponse specifically for unexpected infrastructure failures.
func NewInternalServerError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeInternalServer, message)
}

// NewBadRequestError initializes an ErrorResponse specifically for malformed client requests.
func NewBadRequestError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeBadRequest, message)
}
