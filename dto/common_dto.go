package dto

import "time"

// ===========================================
// STANDARD API RESPONSE STRUCTURE
// ===========================================

/* APIResponse represents the standard API response structure */
type APIResponse struct {
	Success    bool            `json:"success"`
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Data       interface{}     `json:"data,omitempty"`
	Error      *ErrorInfo      `json:"error,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

/* ErrorInfo represents detailed error information */
type ErrorInfo struct {
	Code         string            `json:"code"`
	Message      string            `json:"message"`
	Details      string            `json:"details,omitempty"`
	Validations  []ValidationError `json:"validations,omitempty"`
	RequestID    string            `json:"requestId,omitempty"`
	Timestamp    string            `json:"timestamp"`
}

/* Meta represents metadata for the response (deprecated - use pagination directly) */
type Meta struct {
	RequestID     string          `json:"requestId,omitempty"`
	Timestamp     string          `json:"timestamp"`
	Version       string          `json:"version,omitempty"`
	Pagination    *PaginationMeta `json:"pagination,omitempty"`
	ExecutionTime string          `json:"executionTime,omitempty"`
}

// ===========================================
// RESPONSE BUILDERS
// ===========================================

/* SuccessResponse creates a successful API response */
func SuccessResponse(statusCode int, message string, data interface{}) APIResponse {
	return APIResponse{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

/* SuccessResponseWithPagination creates a successful API response with pagination */
func SuccessResponseWithPagination(statusCode int, message string, data interface{}, pagination *PaginationMeta) APIResponse {
	return APIResponse{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

/* ErrorResponse creates an error API response */
func ErrorResponse(statusCode int, code, message string) APIResponse {
	return APIResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Error: &ErrorInfo{
			Code:      code,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		},
	}
}

/* ErrorResponseWithDetails creates an error API response with details */
func ErrorResponseWithDetails(statusCode int, code, message, details string) APIResponse {
	return APIResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Error: &ErrorInfo{
			Code:      code,
			Message:   message,
			Details:   details,
			Timestamp: getCurrentTimestamp(),
		},
	}
}

/* ValidationErrorResponse creates a validation error response */
func ValidationErrorResponse(validationErrors []ValidationError) APIResponse {
	return APIResponse{
		Success:    false,
		StatusCode: 400,
		Message:    "Validation failed",
		Error: &ErrorInfo{
			Code:        "VALIDATION_ERROR",
			Message:     "Request validation failed",
			Validations: validationErrors,
			Timestamp:   getCurrentTimestamp(),
		},
	}
}

/* NotFoundResponse creates a not found error response */
func NotFoundResponse(resource string) APIResponse {
	return ErrorResponse(404, "NOT_FOUND", resource+" not found")
}

/* UnauthorizedResponse creates an unauthorized error response */
func UnauthorizedResponse() APIResponse {
	return ErrorResponse(401, "UNAUTHORIZED", "Authentication required")
}

/* ForbiddenResponse creates a forbidden error response */
func ForbiddenResponse() APIResponse {
	return ErrorResponse(403, "FORBIDDEN", "Access denied")
}

/* InternalServerErrorResponse creates an internal server error response */
func InternalServerErrorResponse() APIResponse {
	return ErrorResponse(500, "INTERNAL_SERVER_ERROR", "Internal server error occurred")
}

/* ConflictResponse creates a conflict error response */
func ConflictResponse(message string) APIResponse {
	return ErrorResponse(409, "CONFLICT", message)
}

/* BadRequestResponse creates a bad request error response */
func BadRequestResponse(message string) APIResponse {
	return ErrorResponse(400, "BAD_REQUEST", message)
}

// ===========================================
// PAGINATION HELPERS
// ===========================================

/* NewPaginationMeta creates pagination metadata */
func NewPaginationMeta(currentPage, perPage int, totalItems int64) *PaginationMeta {
	totalPages := int((totalItems + int64(perPage) - 1) / int64(perPage))
	
	return &PaginationMeta{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasNextPage: currentPage < totalPages,
		HasPrevPage: currentPage > 1,
	}
}

// ===========================================
// HELPER FUNCTIONS
// ===========================================

/* getCurrentTimestamp returns current timestamp in ISO format */
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// ===========================================
// HTTP STATUS CONSTANTS
// ===========================================

const (
	// Success codes
	StatusOK           = 200
	StatusCreated      = 201
	StatusNoContent    = 204
	
	// Client error codes
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
	
	// Server error codes
	StatusInternalServerError = 500
	StatusBadGateway         = 502
	StatusServiceUnavailable = 503
	StatusGatewayTimeout     = 504
)

// ===========================================
// ERROR CODES
// ===========================================

const (
	// Authentication & Authorization
	ErrorCodeUnauthorized    = "UNAUTHORIZED"
	ErrorCodeForbidden       = "FORBIDDEN"
	ErrorCodeTokenExpired    = "TOKEN_EXPIRED"
	ErrorCodeInvalidToken    = "INVALID_TOKEN"
	
	// Validation
	ErrorCodeValidation      = "VALIDATION_ERROR"
	ErrorCodeBadRequest      = "BAD_REQUEST"
	ErrorCodeInvalidFormat   = "INVALID_FORMAT"
	
	// Resource
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeAlreadyExists   = "ALREADY_EXISTS"
	ErrorCodeConflict        = "CONFLICT"
	
	// Server
	ErrorCodeInternalServer  = "INTERNAL_SERVER_ERROR"
	ErrorCodeDatabaseError   = "DATABASE_ERROR"
	ErrorCodeExternalService = "EXTERNAL_SERVICE_ERROR"
	
	// Rate Limiting
	ErrorCodeRateLimit       = "RATE_LIMIT_EXCEEDED"
	
	// Business Logic
	ErrorCodeBusinessRule    = "BUSINESS_RULE_VIOLATION"
	ErrorCodeInsufficientPermission = "INSUFFICIENT_PERMISSION"
)