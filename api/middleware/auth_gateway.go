package middleware

import (
	"example/api/tool/helper"
	"example/pkg"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(pkg.GetEnv("SECRET")),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "", err.Error())
	}
	return helper.ResponseHandler(c, fiber.StatusUnauthorized, "Invalid or expired JWT", nil)
}
