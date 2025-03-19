package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/dto"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	var userResponse []dto.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	c.IndentedJSON(http.StatusOK, userResponse)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
	}
	userResponse := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.IndentedJSON(http.StatusOK, userResponse)
}
