package main

import (
	"bookkeeper/internal/api"
	"bookkeeper/internal/db"
	"context"
	"log"
)

func main() {

	err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Client.Disconnect(context.TODO())

	r := api.SetUpRouter()

	log.Println("Server running at port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
