package dto

import "time"

// ===========================================
// REQUEST DTOs
// ===========================================

/* CreateUserRequest represents the request structure for creating a user */
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email,max=100"`
	Password  string `json:"password" binding:"required,min=6,max=255"`
	FirstName string `json:"firstName" binding:"max=50"`
	LastName  string `json:"lastName" binding:"max=50"`
}

/* UpdateUserRequest represents the request structure for updating a user */
type UpdateUserRequest struct {
	Username  string `json:"username" binding:"omitempty,min=3,max=50"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
	FirstName string `json:"firstName" binding:"omitempty,max=50"`
	LastName  string `json:"lastName" binding:"omitempty,max=50"`
	IsActive  *bool  `json:"isActive"`
}

/* LoginRequest represents the request structure for user login */
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/* ChangePasswordRequest represents the request structure for changing password */
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=6"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=255"`
}

/* UserSearchRequest represents the request structure for searching users */
type UserSearchRequest struct {
	Query    string `json:"query" form:"query"`
	Page     int    `json:"page" form:"page" binding:"omitempty,min=1"`
	Limit    int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy   string `json:"sortBy" form:"sortBy" binding:"omitempty,oneof=username email firstName lastName isActive createdAt updatedAt"`
	SortDesc bool   `json:"sortDesc" form:"sortDesc"`
	IsActive *bool  `json:"isActive" form:"isActive"`
}

// ===========================================
// RESPONSE DTOs
// ===========================================

/* UserResponse represents the response structure for user data */
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

/* UserListResponse represents the response structure for user list with pagination */
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination PaginationMeta `json:"pagination"`
}

/* UserStatsResponse represents the response structure for user statistics */
type UserStatsResponse struct {
	TotalUsers       int64 `json:"totalUsers"`
	ActiveUsers      int64 `json:"activeUsers"`
	InactiveUsers    int64 `json:"inactiveUsers"`
	DeletedUsers     int64 `json:"deletedUsers"`
	NewUsersLast30d  int64 `json:"newUsersLast30Days"`
	GrowthRate       float64 `json:"growthRate"`
}

/* LoginResponse represents the response structure for user login */
type LoginResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"accessToken"`
	TokenType   string       `json:"tokenType"`
	ExpiresIn   int64        `json:"expiresIn"`
}

// ===========================================
// COMMON DTOs
// ===========================================

/* PaginationMeta represents pagination metadata */
type PaginationMeta struct {
	CurrentPage  int   `json:"currentPage"`
	PerPage      int   `json:"perPage"`
	TotalPages   int   `json:"totalPages"`
	TotalItems   int64 `json:"totalItems"`
	HasNextPage  bool  `json:"hasNextPage"`
	HasPrevPage  bool  `json:"hasPrevPage"`
}

/* ValidationError represents field validation error */
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ===========================================
// VALIDATION HELPERS
// ===========================================

/* Validate validates CreateUserRequest */
func (r *CreateUserRequest) Validate() []ValidationError {
	var errors []ValidationError

	if len(r.Username) < 3 || len(r.Username) > 50 {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "Username must be between 3 and 50 characters",
			Value:   r.Username,
		})
	}

	if len(r.Password) < 6 || len(r.Password) > 255 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be between 6 and 255 characters",
		})
	}

	return errors
}

/* Validate validates UpdateUserRequest */
func (r *UpdateUserRequest) Validate() []ValidationError {
	var errors []ValidationError

	if r.Username != "" && (len(r.Username) < 3 || len(r.Username) > 50) {
		errors = append(errors, ValidationError{
			Field:   "username",
			Message: "Username must be between 3 and 50 characters",
			Value:   r.Username,
		})
	}

	return errors
}

/* SetDefaults sets default values for UserSearchRequest */
func (r *UserSearchRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Limit > 100 {
		r.Limit = 100
	}
	if r.SortBy == "" {
		r.SortBy = "createdAt"
	}
}