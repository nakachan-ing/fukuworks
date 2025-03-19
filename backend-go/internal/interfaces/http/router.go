package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	userRepo := persistence.NewUserRepository(db)
	userHandler := NewUserHandler(userRepo)
	router.GET("/users", userHandler.GetUsers)

	return router

}
