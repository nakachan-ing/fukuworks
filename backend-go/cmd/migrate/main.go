package main

import (
	"fmt"
	"log"

	"github.com/nakachan-ing/fukuworks/backend-go/config"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/infrastructure/persistence"
)

func main() {
	config.LoadConfig()
	fmt.Println("Running in:", config.Env)
	fmt.Println("Database Path:", config.DatabasePath)
	fmt.Println("API URL:", config.ApiUrl)

	_, err := persistence.GetDB(config.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

}
