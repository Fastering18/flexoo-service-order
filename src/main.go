package main

import (
	"log"
	"os"

	"flexoo-service-order/src/routes"
)

func main() {
	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Flexoo Service Order server listening on http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

