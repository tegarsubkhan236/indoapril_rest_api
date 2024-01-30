package routes

import (
	"example/internal/api/handlers"
	"example/internal/api/middleware"
	"example/internal/api/types/permissions"
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
	permissionRoute.Get("/", middleware.Gateway(permissions.READ_PERMISSION), handlers.HandleGetPermissions(crPermissionService))
	permissionRoute.Get("/:id", middleware.Gateway(permissions.READ_PERMISSION), handlers.HandleGetPermission(crPermissionService))
	permissionRoute.Post("/", middleware.Gateway(permissions.CREATE_PERMISSION), handlers.HandleAddPermission(crPermissionService))
	permissionRoute.Put("/:id", middleware.Gateway(permissions.UPDATE_PERMISSION), handlers.HandleUpdatePermission(crPermissionService))
	permissionRoute.Delete("/", middleware.Gateway(permissions.DELETE_PERMISSION), handlers.HandleRemovePermission(crPermissionService))

	roleRoute := api.Group("/v1/core/role", middleware.Protected())
	roleRoute.Get("/", middleware.Gateway(permissions.READ_ROLE), handlers.HandleGetRoles(crRoleService))
	roleRoute.Get("/:id", middleware.Gateway(permissions.READ_ROLE), handlers.HandleGetRole(crRoleService))
	roleRoute.Post("/", middleware.Gateway(permissions.CREATE_ROLE), handlers.HandleAddRole(crRoleService))
	roleRoute.Put("/:id", middleware.Gateway(permissions.UPDATE_ROLE), handlers.HandleUpdateRole(crRoleService))
	roleRoute.Delete("/", middleware.Gateway(permissions.DELETE_ROLE), handlers.HandleRemoveRole(crRoleService))

	userRoute := api.Group("/v1/core/user", middleware.Protected())
	userRoute.Get("/", middleware.Gateway(permissions.READ_USER), handlers.HandleGetUsers(crUserService))
	userRoute.Get("/:id", middleware.Gateway(permissions.READ_USER), handlers.HandleGetUser(crUserService))
	userRoute.Post("/", middleware.Gateway(permissions.CREATE_USER), handlers.HandleAddUser(crUserService))
	userRoute.Put("/:id", middleware.Gateway(permissions.UPDATE_USER), handlers.HandleUpdateUser(crUserService))
	userRoute.Put("/reset_password/:id", middleware.Gateway(permissions.UPDATE_USER), handlers.HandleUpdateUserPassword(crUserService))
	userRoute.Delete("/", middleware.Gateway(permissions.DELETE_USER), handlers.HandleRemoveUser(crUserService))

	teamRoute := api.Group("/v1/core/team", middleware.Protected())
	teamRoute.Get("/", middleware.Gateway(permissions.READ_TEAM), handlers.HandleGetTeams(crTeamService))
	teamRoute.Get("/:id", middleware.Gateway(permissions.READ_TEAM), handlers.HandleGetTeam(crTeamService))
	teamRoute.Post("/", middleware.Gateway(permissions.CREATE_TEAM), handlers.HandleAddTeam(crTeamService))
	teamRoute.Put("/:id", middleware.Gateway(permissions.UPDATE_TEAM), handlers.HandleUpdateTeam(crTeamService))
	teamRoute.Delete("/", middleware.Gateway(permissions.DELETE_TEAM), handlers.HandleRemoveTeam(crTeamService))
}
