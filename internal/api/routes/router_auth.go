package routes

import (
	"example/internal/api/handlers"
	"example/internal/api/middleware"
	"example/internal/pkg/models/cr_user"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, crAuthService cr_user.Service) {
	authRoute := api.Group("/v1/auth")
	authRoute.Post("/login", handlers.Login(crAuthService))
	authRoute.Post("/register", handlers.AddUser(crAuthService))

	profileRoute := api.Group("/v1/profile", middleware.Protected())
	profileRoute.Get("/", handlers.Me(crAuthService))
	profileRoute.Put("/update_info", handlers.UpdateMe(crAuthService))
	profileRoute.Put("/update_password", handlers.UpdatePasswordMe(crAuthService))
}
