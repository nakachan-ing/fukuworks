package main

import (
	"fmt"
	"log"

	"github.com/nakachan-ing/fukuworks/backend-go/config"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/infrastructure/persistence"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/interfaces/http"
)

func main() {
	config.LoadConfig()

	db, err := persistence.GetDB(config.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := http.NewRouter(db)
	apiUrl := config.ApiUrl
	fmt.Println("Server is running on  " + apiUrl)
	router.Run(apiUrl)

}
