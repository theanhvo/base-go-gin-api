package handlers

import (
	"strconv"

	"baseApi/dto"
	"baseApi/logger"
	"baseApi/middleware"
	"baseApi/monitoring"
	"baseApi/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

/* NewUserHandler creates a new user handler */
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

/* CreateUser handles user creation */
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body:", err)
		response := dto.ValidationErrorResponse([]dto.ValidationError{
			{Field: "request", Message: "Invalid request format", Value: err.Error()},
		})
		c.JSON(response.StatusCode, response)
		return
	}

	// Additional validation
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		response := dto.ValidationErrorResponse(validationErrors)
		c.JSON(response.StatusCode, response)
		return
	}

	// Start Sentry span for service call
	span := middleware.StartSpanFromContext(c, "user.create", "Create new user")
	user, err := h.userService.CreateUser(req)
	if span != nil {
		span.Finish()
	}

	if err != nil {
		// Capture error to Sentry with context
		monitoring.CaptureError(err, map[string]interface{}{
			"operation": "create_user",
			"username":  req.Username,
			"email":     req.Email,
		})

		logger.Error("Failed to create user:", err)
		response := dto.ErrorResponseWithDetails(
			dto.StatusInternalServerError,
			dto.ErrorCodeDatabaseError,
			"Failed to create user",
			err.Error(),
		)
		c.JSON(response.StatusCode, response)
		return
	}

	logger.Info("User created successfully:", user.ID)
	response := dto.SuccessResponse(
		dto.StatusCreated,
		"User created successfully",
		user,
	)
	c.JSON(response.StatusCode, response)
}

/* GetUser handles retrieving a user by ID */
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response := dto.BadRequestResponse("Invalid user ID format")
		c.JSON(response.StatusCode, response)
		return
	}

	// Start Sentry span for service call
	span := middleware.StartSpanFromContext(c, "user.get_by_id", "Get user by ID")
	user, err := h.userService.GetUserByID(uint(id))
	if span != nil {
		span.Finish()
	}

	if err != nil {
		if err.Error() == "user not found" {
			// Add breadcrumb for not found
			monitoring.AddBreadcrumb("User not found", "user", map[string]interface{}{
				"user_id": id,
			})
			response := dto.NotFoundResponse("User")
			c.JSON(response.StatusCode, response)
			return
		}

		// Capture error to Sentry
		monitoring.CaptureError(err, map[string]interface{}{
			"operation": "get_user_by_id",
			"user_id":   id,
		})

		logger.Error("Failed to get user:", err)
		response := dto.ErrorResponseWithDetails(
			dto.StatusInternalServerError,
			dto.ErrorCodeDatabaseError,
			"Failed to retrieve user",
			err.Error(),
		)
		c.JSON(response.StatusCode, response)
		return
	}

	response := dto.SuccessResponse(dto.StatusOK, "User retrieved successfully", user)
	c.JSON(response.StatusCode, response)
}

/* GetUserByUsername handles retrieving a user by username */
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		response := dto.BadRequestResponse("Username parameter is required")
		c.JSON(response.StatusCode, response)
		return
	}

	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "user not found" {
			response := dto.NotFoundResponse("User")
			c.JSON(response.StatusCode, response)
			return
		}
		logger.Error("Failed to get user:", err)
		response := dto.ErrorResponseWithDetails(
			dto.StatusInternalServerError,
			dto.ErrorCodeDatabaseError,
			"Failed to retrieve user",
			err.Error(),
		)
		c.JSON(response.StatusCode, response)
		return
	}

	response := dto.SuccessResponse(dto.StatusOK, "User retrieved successfully", user)
	c.JSON(response.StatusCode, response)
}

/* GetAllUsers handles retrieving all users with pagination and search */
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var searchReq dto.UserSearchRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&searchReq); err != nil {
		response := dto.BadRequestResponse("Invalid query parameters")
		c.JSON(response.StatusCode, response)
		return
	}

	// Set defaults and validate
	searchReq.SetDefaults()

	// Start Sentry span for service call
	span := middleware.StartSpanFromContext(c, "user.get_all", "Get all users with search")
	userList, err := h.userService.GetAllUsers(searchReq)
	if span != nil {
		span.Finish()
	}

	if err != nil {
		// Capture error to Sentry with search context
		monitoring.CaptureError(err, map[string]interface{}{
			"operation":   "get_all_users",
			"search_query": searchReq.Query,
			"page":        searchReq.Page,
			"limit":       searchReq.Limit,
		})

		logger.Error("Failed to get users:", err)
		response := dto.ErrorResponseWithDetails(
			dto.StatusInternalServerError,
			dto.ErrorCodeDatabaseError,
			"Failed to retrieve users",
			err.Error(),
		)
		c.JSON(response.StatusCode, response)
		return
	}

	response := dto.SuccessResponseWithPagination(
		dto.StatusOK,
		"Users retrieved successfully",
		userList.Users,
		&userList.Pagination,
	)
	c.JSON(response.StatusCode, response)
}

/* UpdateUser handles user updates */
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response := dto.BadRequestResponse("Invalid user ID format")
		c.JSON(response.StatusCode, response)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body:", err)
		response := dto.ValidationErrorResponse([]dto.ValidationError{
			{Field: "request", Message: "Invalid request format", Value: err.Error()},
		})
		c.JSON(response.StatusCode, response)
		return
	}

	// Additional validation
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		response := dto.ValidationErrorResponse(validationErrors)
		c.JSON(response.StatusCode, response)
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		if err.Error() == "user not found" {
			response := dto.NotFoundResponse("User")
			c.JSON(response.StatusCode, response)
			return
		}
		logger.Error("Failed to update user:", err)
		response := dto.ErrorResponseWithDetails(
			dto.StatusInternalServerError,
			dto.ErrorCodeDatabaseError,
			"Failed to update user",
			err.Error(),
		)
		c.JSON(response.StatusCode, response)
		return
	}

	logger.Info("User updated successfully:", user.ID)
	response := dto.SuccessResponse(
		dto.StatusOK,
		"User updated successfully",
		user,
	)
	c.JSON(response.StatusCode, response)
}

/* DeleteUser handles user deletion */
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response := dto.BadRequestResponse("Invalid user ID format")
		c.JSON(response.StatusCode, response)
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			response := dto.NotFoundResponse("User")
			c.JSON(response.StatusCode, response)
			return
		}
		logger.Error("Failed to delete user:", err)
		response := dto.ErrorResponseWithDetails(
			dto.StatusInternalServerError,
			dto.ErrorCodeDatabaseError,
			"Failed to delete user",
			err.Error(),
		)
		c.JSON(response.StatusCode, response)
		return
	}

	logger.Info("User deleted successfully:", id)
	response := dto.SuccessResponse(
		dto.StatusOK,
		"User deleted successfully",
		nil,
	)
	c.JSON(response.StatusCode, response)
}