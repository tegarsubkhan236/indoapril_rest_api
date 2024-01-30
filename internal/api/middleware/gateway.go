package middleware

import (
	"example/internal/api/types/roles"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Gateway(permission string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error - Invalid Token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error - Invalid Claims"})
		}

		role, ok := claims["role_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error - Invalid Role ID"})
		}

		// SUPER_ADMIN
		if role == roles.SUPER_ADMIN {
			return c.Next()
		}

		permissions, ok := claims["permissions"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error - Invalid Permissions"})
		}

		// Has Permission
		for _, p := range permissions {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden - You have no access"})
	}
}
