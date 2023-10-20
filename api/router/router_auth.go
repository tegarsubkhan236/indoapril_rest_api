package router

import (
	"example/api/controller"
	"github.com/gofiber/fiber/v2"
)

func setupAuthRoutes(v1 fiber.Router) {
	authRoute := v1.Group("/auth")
	authRoute.Post("/login", controller.Login)
	authRoute.Post("/register", controller.CreateUser)
}
