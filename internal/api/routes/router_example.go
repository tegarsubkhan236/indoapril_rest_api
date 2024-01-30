package routes

import (
	"example/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupExampleRoutes(db *gorm.DB, api fiber.Router) {
	routeEx := api.Group("/v1/example")
	routeEx.Get("/get-so", handlers.HandleIzziGetSO(db))
}
