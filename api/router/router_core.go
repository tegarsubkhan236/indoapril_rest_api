package router

import (
	"example/api/controller"
	"example/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupCoreRoutes(v1 fiber.Router) {
	coreRoute := v1.Group("/core", middleware.Protected())

	userRoute := coreRoute.Group("/user")
	userRoute.Get("/", controller.GetUsers)
	userRoute.Get("/:id", controller.GetUser)
	userRoute.Post("/", controller.CreateUser)
	userRoute.Put("/:id", controller.UpdateUser)
	userRoute.Put("/reset_password/:id", controller.UpdateUserPassword)
	userRoute.Delete("/", controller.DeleteUser)

	permissionRoute := coreRoute.Group("/permission")
	permissionRoute.Get("/", controller.GetPermissions)
	permissionRoute.Get("/:id", controller.GetPermission)
	permissionRoute.Post("/", controller.CreatePermission)
	permissionRoute.Put("/:id", controller.UpdatePermission)
	permissionRoute.Delete("/", controller.DeletePermission)

	roleRoute := coreRoute.Group("/role")
	roleRoute.Get("/", controller.GetRoles)
	roleRoute.Get("/:id", controller.GetRole)
	roleRoute.Post("/", controller.CreateRole)
	roleRoute.Put("/:id", controller.UpdateRole)
	roleRoute.Delete("/", controller.DeleteRole)

	//teamRoute := coreRoute.Group("/team")
	//teamRoute.Get("/", core.GetTeams)
	//teamRoute.Get("/:id", core.GetTeam)
	//teamRoute.Post("/", core.CreateTeam)
	//teamRoute.Put("/:id", core.UpdateTeam)
	//teamRoute.Delete("/", core.DeleteTeam)
}
