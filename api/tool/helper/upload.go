package helper

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

func SingleUpload(c *fiber.Ctx, paramName, oldFilePath string) (string, error) {
	file, errParseFile := c.FormFile(paramName)
	if errParseFile != nil {
		return "", errParseFile
	}

	if file != nil {
		uploadDir := "F:/www/indoapril/upload"
		fileNewPath := filepath.Join(uploadDir, file.Filename)

		if errSaveFile := c.SaveFile(file, fileNewPath); errSaveFile != nil {
			return "", errSaveFile
		}

		if oldFilePath != "" {
			if err := os.Remove(oldFilePath); err != nil {
				return "", err
			}
		}

		return fileNewPath, nil
	}

	return "", nil
}
