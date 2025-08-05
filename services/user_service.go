package services

import (
	"errors"
	"fmt"
	"time"

	"baseApi/cache"
	"baseApi/database"
	"baseApi/dto"
	"baseApi/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

/* NewUserService creates a new user service instance */
func NewUserService() *UserService {
	return &UserService{}
}

/* CreateUser creates a new user */
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	user.FromCreateDTO(req)
	user.Password = string(hashedPassword) // Override with hashed password

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// Cache user data
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	cache.Set(cacheKey, user, 1*time.Hour)

	response := user.ToDTO()
	return &response, nil
}

/* GetUserByID retrieves a user by ID with caching */
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d", id)
	var cachedUser models.User
	if err := cache.Get(cacheKey, &cachedUser); err == nil {
		response := cachedUser.ToDTO()
		return &response, nil
	}

	// If not in cache, get from database
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Cache the user
	cache.Set(cacheKey, user, 1*time.Hour)

	response := user.ToDTO()
	return &response, nil
}

/* GetUserByUsername retrieves a user by username */
func (s *UserService) GetUserByUsername(username string) (*dto.UserResponse, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := user.ToDTO()
	return &response, nil
}

/* GetAllUsers retrieves all users with pagination */
func (s *UserService) GetAllUsers(req dto.UserSearchRequest) (*dto.UserListResponse, error) {
	req.SetDefaults()

	var users []models.User
	var totalCount int64

	query := database.DB.Model(&models.User{})

	// Apply search filter
	if req.Query != "" {
		searchTerm := "%" + req.Query + "%"
		query = query.Where(
			"username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Apply active filter
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting with field mapping
	var orderField string
	switch req.SortBy {
	case "createdAt":
		orderField = "created_at"
	case "updatedAt":
		orderField = "updated_at"
	case "firstName":
		orderField = "first_name"
	case "lastName":
		orderField = "last_name"
	case "isActive":
		orderField = "is_active"
	default:
		orderField = req.SortBy // username, email use same name
	}
	
	orderClause := orderField
	if req.SortDesc {
		orderClause += " DESC"
	} else {
		orderClause += " ASC"
	}
	query = query.Order(orderClause)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&users).Error; err != nil {
		return nil, err
	}

	// Create pagination metadata
	pagination := dto.NewPaginationMeta(req.Page, req.Limit, totalCount)

	// Convert to DTO
	response := models.ToUserListDTO(users, pagination)
	return &response, nil
}

/* UpdateUser updates a user */
func (s *UserService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields using DTO
	user.UpdateFromDTO(req)

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	// Update cache
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	cache.Set(cacheKey, user, 1*time.Hour)

	response := user.ToDTO()
	return &response, nil
}

/* DeleteUser soft deletes a user */
func (s *UserService) DeleteUser(id uint) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return err
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("user:%d", id)
	cache.Delete(cacheKey)

	return nil
}

/* GetUserCount returns the total number of users */
func (s *UserService) GetUserCount() (int64, error) {
	var count int64
	if err := database.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}