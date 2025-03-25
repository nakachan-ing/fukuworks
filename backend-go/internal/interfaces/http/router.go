package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/infrastructure/persistence"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http/middleware"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ReservedPathGuard())

	userRepo := persistence.NewUserRepository(db)
	userHandler := NewUserHandler(userRepo)

	projectRepo := persistence.NewProjectRepository(db)
	projectHandler := NewProjectHandler(projectRepo)

	taskRepo := persistence.NewTaskRepository(db)
	taskHandler := NewTaskHandler(taskRepo)

	// ==============================================================================
	// for user

	// 認証なしで使えるパス（公開）
	router.POST("/signup", userHandler.PostUser)
	router.POST("/login", userHandler.Login)

	// 認証が必要なパス
	authorized := router.Group("/:user")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("", userHandler.GetUser)
		authorized.PATCH("", userHandler.UpdateUser)
		authorized.DELETE("", userHandler.SoftDeleteUser)

		authorized.POST("/projects", projectHandler.PostProject)
		authorized.GET("/projects", projectHandler.GetAllProjectsByUser)
		authorized.GET("/projects/:pid", projectHandler.GetProject)
		authorized.PATCH("/projects/:pid", projectHandler.UpdateProject)
		authorized.DELETE("/projects/:pid", projectHandler.SoftDeleteProject)

		authorized.POST("/projects/:pid/tasks", taskHandler.PostTask)
		authorized.GET("/projects/:pid/tasks", taskHandler.GetAllTasksByProject)
		authorized.GET("/projects/:pid/tasks/:tid", taskHandler.GetTask)
		authorized.PATCH("/projects/:pid/tasks/:tid", taskHandler.UpdateTask)
		authorized.DELETE("/projects/:pid/tasks/:tid", taskHandler.SoftDeleteTask)
	}

	// ==============================================================================

	// ==============================================================================
	// for owner
	api := router.Group("/admin")
	{
		api.GET("/users", userHandler.GetAllUsers)
		api.DELETE("/users/:id", userHandler.HardDeleteUser)
		api.GET("/projects", projectHandler.GetAllProjectsForOwner)
		api.DELETE("/projects/:id", projectHandler.HardDeleteProject)
		api.GET("/tasks", taskHandler.GetAllTasksForOwner)
		api.DELETE("/tasks/:id", taskHandler.HardDeleteTask)
	}
	// ==============================================================================

	return router

}
