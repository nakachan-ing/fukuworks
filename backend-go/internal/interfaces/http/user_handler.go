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
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
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
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	if err := h.userRepo.Create(&newUser); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
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

func (h *UserHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequest
	if !userBindAndValidate(c, &loginReq) {
		return
	}

	user, err := h.userRepo.Find(loginReq.Name)
	if err != nil || user.Password != loginReq.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 仮のトークンを返す（JWT対応は今後）
	token := "mock-token-for-" + user.Name
	c.JSON(http.StatusOK, gin.H{"token": token})
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userName := c.Param("user")
	var userRequest dto.UserUpdateRequest
	if !userBindAndValidate(c, &userRequest) {
		return
	}

	targetUser := models.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	updatedUser, err := h.userRepo.Update(userName, &targetUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
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

func (h *UserHandler) SoftDeleteUser(c *gin.Context) {
	userName := c.Param("user")
	if err := h.userRepo.SoftDelete(userName); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}

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
