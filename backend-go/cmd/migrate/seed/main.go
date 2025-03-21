package main

import (
	"fmt"
	"log"

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

	log.Println("Seeding completed!")

}
