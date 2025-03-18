package main

import (
	"fmt"

	"github.com/nakachan-ing/fukuworks/backend-go/config"
)

func main() {
	config.LoadConfig()
	fmt.Println("Running in:", config.Env)
	fmt.Println("Database Path:", config.DatabasePath)
	fmt.Println("API URL:", config.ApiUrl)

	// db, err := persistence.GetDB(config.DatabasePath)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

}
