package main

import (
	"febre/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Start a new fiber app
	app := fiber.New()

	// Connect to the Database
	//database.ConnectDB()

	//Connect to the Redis

	// Setup the router
	router.SetupRoutes(app)

	// Listen on PORT 300
	app.Listen(":3000")
}
