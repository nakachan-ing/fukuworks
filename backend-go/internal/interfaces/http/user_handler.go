package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// for user
func (h *UserHandler) PostUser(c *gin.Context) {
	var userRequest dto.UserCreateRequest
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	newUser := models.User{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	}

	err := h.userRepo.Create(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	userResponse := dto.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userName := c.Param("user")
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	// }

	user, err := h.userRepo.Find(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
	}
	userResponse := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.IndentedJSON(http.StatusOK, userResponse)
}

// for user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userName := c.Param("user")
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	// }

	var userRequest dto.UserUpdateRequest
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	targetUser := models.User{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	}

	updatedUser, err := h.userRepo.Update(userName, &targetUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	userResponse := dto.UserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

// for user
func (h *UserHandler) SoftDeleteUser(c *gin.Context) {
	userName := c.Param("user")
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	// }

	err := h.userRepo.SoftDelete(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
	}

	c.Status(http.StatusNoContent)
}

// for owner
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	var userResponse []dto.UserOwnerResponse
	for _, user := range users {
		userResponse = append(userResponse, dto.UserOwnerResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	c.IndentedJSON(http.StatusOK, userResponse)
}

func (h *UserHandler) HardDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	err = h.userRepo.HardDelete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
	}

	c.Status(http.StatusNoContent)
}
