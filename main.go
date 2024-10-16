package main

import (
	"kenshilabs/database"
	"kenshilabs/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	client := database.ConnectDB()
	if client == nil {
		log.Fatal("Failed to connect to MongoDB")
	}

	app := fiber.New()

	routes.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
