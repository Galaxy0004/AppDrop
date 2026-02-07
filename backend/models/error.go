package models

// ErrorCode represents the type of error
type ErrorCode string

const (
	ErrorCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeConflict       ErrorCode = "CONFLICT"
	ErrorCodeInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest     ErrorCode = "BAD_REQUEST"
)

// APIError represents the error response format
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// ErrorResponse represents the API error response structure
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	}
}

// NewValidationError creates a validation error response
func NewValidationError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeValidation, message)
}

// NewNotFoundError creates a not found error response
func NewNotFoundError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeNotFound, message)
}

// NewConflictError creates a conflict error response
func NewConflictError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeConflict, message)
}

// NewInternalServerError creates an internal server error response
func NewInternalServerError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeInternalServer, message)
}

// NewBadRequestError creates a bad request error response
func NewBadRequestError(message string) ErrorResponse {
	return NewErrorResponse(ErrorCodeBadRequest, message)
}
