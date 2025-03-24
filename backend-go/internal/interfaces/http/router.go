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

	taskRepo := persistence.NewTaskRepository(db)
	taskHandler := NewTaskHandler(taskRepo)

	// for user
	router.POST("/users", userHandler.PostUser)
	router.GET("/:user", userHandler.GetUser)
	router.PATCH("/:user", userHandler.UpdateUser)
	router.DELETE("/:user", userHandler.SoftDeleteUser)

	router.POST("/:user/projects", projectHandler.PostProject)
	router.GET("/:user/projects", projectHandler.GetAllProjectsByUser)
	router.GET("/:user/projects/:id", projectHandler.GetProject)
	router.PATCH("/:user/:projects/:id", projectHandler.UpdateProject)
	router.DELETE("/:user/:project/:id", projectHandler.SoftDeleteProject)

	router.POST("/:user/:projects/:id/tasks", taskHandler.PostTask)
	router.GET("/:user/:projects/:pid/tasks", taskHandler.GetAllTasksByProject)
	router.GET("/:user/:projects/:pid/tasks/:tid", taskHandler.GetTask)
	router.PATCH("/user/:projects/:pid/tasks/:tid,", taskHandler.UpdateTask)
	router.DELETE("/:user/:projects/:pid/tasks/:tid", taskHandler.SoftDeleteTask)

	// for owner
	api := router.Group("/api")
	{
		api.GET("/users", userHandler.GetAllUsers)
		api.DELETE("/users/:id", userHandler.HardDeleteUser)
		api.GET("/projects", projectHandler.GetAllProjectsForOwner)
		api.DELETE("/projects/:id", projectHandler.HardDeleteProject)
		api.GET("/tasks", taskHandler.GetAllTasksForOwner)
		api.DELETE("/tasks/:id")
	}
	return router

}
