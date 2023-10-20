package main

import (
	"example/api/router"
	"example/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Authorization, Content-Type",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))
	pkg.ConnectDB()

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":4000"))
}
