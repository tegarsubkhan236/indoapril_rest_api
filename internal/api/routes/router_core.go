package routes

import (
	"example/internal/api/handlers"
	"example/internal/api/middleware"
	"example/internal/api/util/constant"
	"example/internal/pkg/models/cr_permission"
	"example/internal/pkg/models/cr_role"
	"example/internal/pkg/models/cr_team"
	"example/internal/pkg/models/cr_user"
	"github.com/gofiber/fiber/v2"
)

func SetupCoreRoutes(
	api fiber.Router,
	crPermissionService cr_permission.Service,
	crRoleService cr_role.Service,
	crUserService cr_user.Service,
	crTeamService cr_team.Service,
) {
	permissionRoute := api.Group("/v1/core/permission", middleware.Protected())
	permissionRoute.Get("/", handlers.GetPermissions(crPermissionService))
	permissionRoute.Get("/:id", handlers.GetPermission(crPermissionService))
	permissionRoute.Post("/", handlers.AddPermission(crPermissionService))
	permissionRoute.Put("/:id", handlers.UpdatePermission(crPermissionService))
	permissionRoute.Delete("/", handlers.RemovePermission(crPermissionService))

	roleRoute := api.Group("/v1/core/role", middleware.Protected())
	roleRoute.Get("/", handlers.GetRoles(crRoleService))
	roleRoute.Get("/:id", handlers.GetRole(crRoleService))
	roleRoute.Post("/", handlers.AddRole(crRoleService))
	roleRoute.Put("/:id", handlers.UpdateRole(crRoleService))
	roleRoute.Delete("/", handlers.RemoveRole(crRoleService))

	userRoute := api.Group("/v1/core/user", middleware.Protected())
	userRoute.Get("/", middleware.Gateway(constant.READ_USER), handlers.GetUsers(crUserService))
	userRoute.Get("/:id", middleware.Gateway(constant.READ_USER), handlers.GetUser(crUserService))
	userRoute.Post("/", middleware.Gateway(constant.CREATE_USER), handlers.AddUser(crUserService))
	userRoute.Put("/:id", middleware.Gateway(constant.UPDATE_USER), handlers.UpdateUser(crUserService))
	userRoute.Put("/reset_password/:id", middleware.Gateway(constant.UPDATE_USER), handlers.UpdateUserPassword(crUserService))
	userRoute.Delete("/", middleware.Gateway(constant.DELETE_USER), handlers.RemoveUser(crUserService))

	//teamRoute := api.Group("/v1/core/team", middleware.Protected())
	//teamRoute.Get("/", handlers.GetTeams(crTeamService))
	//teamRoute.Get("/:id", handlers.GetTeam(crTeamService))
	//teamRoute.Post("/", handlers.AddTeam(crTeamService))
	//teamRoute.Put("/:id", handlers.Update(crTeamService))
	//teamRoute.Delete("/", handlers.RemoveTeam(crTeamService))
}
