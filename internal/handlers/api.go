package handlers

import (
	"net/http"
	"strconv"

	"{{MODULE_NAME}}/internal/models"

	"github.com/velocitykode/velocity/pkg/auth"
	"github.com/velocitykode/velocity/pkg/log"
	"github.com/velocitykode/velocity/pkg/router"
)

// Health returns the API health status
func Health(ctx *router.Context) error {
	return ctx.JSON(router.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// ListUsers returns all users
func ListUsers(ctx *router.Context) error {
	users, err := models.User{}.All()
	if err != nil {
		log.Error("Failed to fetch users", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch users",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": users,
	})
}

// GetUser returns a single user by ID
func GetUser(ctx *router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	user, err := models.User{}.Find(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}

// CreateUser creates a new user
func CreateUser(ctx *router.Context) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": "Name, email, and password are required",
		})
	}

	// Check if user already exists
	existingUser, _ := models.User{}.FindBy("email", input.Email)
	if existingUser != nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"error": "A user with this email already exists",
		})
	}

	// Hash password
	hashedPassword, err := auth.Hash(input.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process password",
		})
	}

	// Create user
	user, err := models.User{}.Create(map[string]any{
		"name":     input.Name,
		"email":    input.Email,
		"password": hashedPassword,
	})
	if err != nil {
		log.Error("Failed to create user", "error", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"data": user,
	})
}

// GetCurrentUser returns the authenticated user
func GetCurrentUser(ctx *router.Context) error {
	user := auth.User(ctx.Request)
	if user == nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Not authenticated",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": user,
	})
}
