package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nakachan-ing/fukuworks/backend-go/config"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/infrastructure/persistence"
)

func main() {
	config.LoadConfig()
	fmt.Println("Running in:", config.Env)
	fmt.Println("Database Path:", config.DatabasePath)
	fmt.Println("API URL:", config.ApiUrl)

	db, err := persistence.GetDB(config.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Seeding database...")

	userRepo := persistence.NewUserRepository(db)
	users := []models.User{
		{Name: "nakachan-ing", Email: "hogehoge@gmail.com"},
		{Name: "nakachan", Email: "hogehoge.icloud.com"},
	}
	for _, user := range users {
		userRepo.Create(&user)
	}

	projectRepo := persistence.NewProjectRepository(db)
	projects := []models.Project{
		{
			UserID:       1,
			Title:        "FukuWarksの開発",
			Description:  "副業案件管理アプリケーションの作成プロジェクト",
			Platform:     "個人",
			Client:       "個人",
			EstimatedFee: 0,
			Status:       "In progress",
			Deadline:     time.Date(2025, 5, 31, 0, 0, 0, 0, time.Local),
		},
		{
			UserID:       2,
			Title:        "ztl-cliの改修",
			Description:  "zettelkasten-cliのプロジェクト",
			Platform:     "個人",
			Client:       "個人",
			EstimatedFee: 0,
			Status:       "Open",
			Deadline:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
		},
	}
	projectRepo.Create("nakachan-ing", &projects[0])
	projectRepo.Create("nakachan", &projects[1])

	taskRepo := persistence.NewTaskRepository(db)
	tasks := []models.Task{
		{
			Title:       "Routingの実装",
			Description: "APIを叩くためのルーティングの実装",
			Status:      "Todo",
			Priority:    "High",
			DueDate:     time.Date(2025, 3, 25, 0, 0, 0, 0, time.Local),
		},
		{
			Title:       "ztl syncコマンドの改修",
			Description: "S3にアップロード、ダウンロード機能の改修",
			Status:      "Todo",
			Priority:    "Medium",
			DueDate:     time.Date(2025, 4, 30, 0, 0, 0, 0, time.Local),
		},
	}
	taskRepo.Create("nakachan-ing", 1, &tasks[0])
	taskRepo.Create("nakachan", 1, &tasks[1])

	log.Println("Seeding completed!")

}
