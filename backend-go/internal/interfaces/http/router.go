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

	projectRepo := persistence.NewProjectRepository(db)
	projectHandler := NewProjectHandler(projectRepo)

	// for user
	router.POST("/users", userHandler.PostUser)         // for user
	router.GET("/:user", userHandler.GetUser)           // for user
	router.PATCH("/:user", userHandler.UpdateUser)      // for user
	router.DELETE("/:user", userHandler.SoftDeleteUser) // for user

	router.POST("/projects", projectHandler.PostProject)
	router.GET("/:user/:project", projectHandler.GetProject) // for user

	// for owner
	api := router.Group("/api")
	{
		api.GET("/users", userHandler.GetAllUsers)           // for owner
		api.DELETE("/users/:id", userHandler.HardDeleteUser) // for owner
	}
	return router

}
