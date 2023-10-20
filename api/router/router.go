package router

import (
	"example/api/tool/helper"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"path/filepath"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	v1.Post("/", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		files := form.File["documents"]
		uploadDir := "F:/www/indoapril/upload"
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			err := c.SaveFile(file, filepath.Join(uploadDir, file.Filename))
			if err != nil {
				return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err.Error())
			}
		}
		return helper.ResponseHandler(c, fiber.StatusOK, "Upload Success", nil)
	})

	setupAuthRoutes(v1)
	setupCoreRoutes(v1)
	setupMasterRoutes(v1)
	setupTransactionRoutes(v1)
}
