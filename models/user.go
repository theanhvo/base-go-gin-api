package models

import (
	"time"

	"baseApi/dto"

	"gorm.io/gorm"
)

/* User represents the user model in the database */
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null;size:50"`
	Email     string         `json:"email" gorm:"unique;not null;size:100"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	FirstName string         `json:"firstName" gorm:"column:first_name;size:50"`
	LastName  string         `json:"lastName" gorm:"column:last_name;size:50"`
	IsActive  bool           `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

/* TableName specifies the table name for User model */
func (User) TableName() string {
	return "users"
}

/* ToDTO converts User model to UserResponse DTO */
func (u *User) ToDTO() dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

/* FromCreateDTO creates User model from CreateUserRequest DTO */
func (u *User) FromCreateDTO(req dto.CreateUserRequest) {
	u.Username = req.Username
	u.Email = req.Email
	u.Password = req.Password // Note: Password should be hashed before calling this
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.IsActive = true
}

/* UpdateFromDTO updates User model from UpdateUserRequest DTO */
func (u *User) UpdateFromDTO(req dto.UpdateUserRequest) {
	if req.Username != "" {
		u.Username = req.Username
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
}

/* ToListDTO converts slice of User models to UserListResponse DTO */
func ToUserListDTO(users []User, pagination *dto.PaginationMeta) dto.UserListResponse {
	userDTOs := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToDTO()
	}

	return dto.UserListResponse{
		Users:      userDTOs,
		Pagination: *pagination,
	}
}