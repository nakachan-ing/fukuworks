package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/infrastructure/persistence"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http/middleware"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://host.docker.internal:8080"}, // Appsmith URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// router.Use(middleware.ReservedPathGuard())

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
	authorized.Use(middleware.ReservedPathGuard(), middleware.AuthMiddleware())
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
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireAdminMiddleware())
	{
		admin.GET("/me", userHandler.GetAdminMe)
		admin.GET("/users", userHandler.GetAllUsers)
		admin.DELETE("/users/:id", userHandler.HardDeleteUser)
		admin.GET("/projects", projectHandler.GetAllProjectsForOwner)
		admin.DELETE("/projects/:id", projectHandler.HardDeleteProject)
		admin.GET("/tasks", taskHandler.GetAllTasksForOwner)
		admin.DELETE("/tasks/:id", taskHandler.HardDeleteTask)
	}
	// ==============================================================================

	return router

}
