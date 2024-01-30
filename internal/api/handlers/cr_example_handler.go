package handlers

import (
	"example/internal/api/util/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type trSoH struct {
	ID        uint      `json:"id"`
	SoCode    string    `json:"so_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated-at"`
}

func HandleIzziGetSO(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, count, err := func() ([]trSoH, int64, error) {
			var data []trSoH

			db = db.Table("tr_so_h")

			if err := db.Find(&data).Error; err != nil {
				return nil, 0, err
			}

			return data, 0, nil
		}()

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(response.ErrorResponse(err))
		}

		c.Status(http.StatusOK)
		return c.JSON(response.SuccessesResponse(result, count))
	}
}
