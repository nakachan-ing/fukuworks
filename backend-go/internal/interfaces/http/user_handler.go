package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/dto"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// Reserved words (e.g., for static routes like /login, /admin, etc.)
var reservedPaths = map[string]bool{
	"login":  true,
	"admin":  true,
	"health": true,
}

// Middleware to skip reserved paths from user route handling
func ReservedPathGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		first := strings.Split(strings.TrimLeft(c.Request.URL.Path, "/"), "/")[0]
		if reservedPaths[first] {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Route not found"})
			return
		}
		c.Next()
	}
}

func userBindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errors := make([]map[string]string, len(ve))
			for i, fe := range ve {
				errors[i] = map[string]string{
					"field":   fe.Field(),
					"message": validationErrorMessage(fe),
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return false
	}
	return true
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	default:
		return "is invalid"
	}
}

// ==================================================================================================================
// for user
func (h *UserHandler) PostUser(c *gin.Context) {
	var userRequest dto.UserCreateRequest
	if !userBindAndValidate(c, &userRequest) {
		return
	}

	newUser := models.User{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	}

	if err := h.userRepo.Create(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userResponse := dto.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: newUser.UpdatedAt.Format(time.RFC3339),
	}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userName := c.Param("user")
	user, err := h.userRepo.Find(userName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	c.IndentedJSON(http.StatusOK, userResponse)
}

// for user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userName := c.Param("user")
	var userRequest dto.UserUpdateRequest
	if !userBindAndValidate(c, &userRequest) {
		return
	}

	targetUser := models.User{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	}

	updatedUser, err := h.userRepo.Update(userName, &targetUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	userResponse := dto.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
	}

	c.IndentedJSON(http.StatusOK, userResponse)
}

// for user
func (h *UserHandler) SoftDeleteUser(c *gin.Context) {
	userName := c.Param("user")
	if err := h.userRepo.SoftDelete(userName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ==================================================================================================================

// ==================================================================================================================
// for owner
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	var userResponse []dto.UserResponseForOwner
	for _, user := range users {
		userResponse = append(userResponse, dto.UserResponseForOwner{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
			DeletedAt: user.DeletedAt.Time.Format(time.RFC3339),
		})
	}
	c.IndentedJSON(http.StatusOK, userResponse)
}

func (h *UserHandler) HardDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
		return
	}

	if err := h.userRepo.HardDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
