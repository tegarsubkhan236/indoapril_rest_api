package temp

import (
	"example/api/tool/helper"
	"github.com/gofiber/fiber/v2"
)

func BalanceCheck(c *fiber.Ctx) error {
	body, err := HitBalanceSCheck()
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "Error on server", err.Error())
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "", body.Data)
}

func PriceList(c *fiber.Ctx) error {
	var code = c.Query("code", "")
	var priceType = c.Query("type", "")
	var groupBy = c.Query("groupBy", "")

	body, err := HitPriceList(priceType, code)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "Error on server", err.Error())
	}

	if groupBy != "" {
		body, err := PriceListGrouped(body.Data, groupBy)
		if err != nil {
			return helper.ResponseHandler(c, fiber.StatusBadRequest, "Bad request", err.Error())
		} else {
			return helper.ResponseHandler(c, fiber.StatusOK, "", body)
		}
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "", body.Data)
}

func Deposit(c *fiber.Ctx) error {
	body, err := HitDeposit()
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "Error on server", err.Error())
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "", body.Data)
}
