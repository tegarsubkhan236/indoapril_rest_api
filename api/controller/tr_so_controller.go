package controller

import (
	"example/api/service"
	"example/api/tool/converter"
	"example/api/tool/helper"
	"example/pkg/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func GetSalesOrders(c *fiber.Ctx) error {
	var req struct {
		Page         int                `query:"page" default:"1"`
		Limit        int                `query:"limit" default:"10"`
		TrSalesOrder model.TrSalesOrder `query:"column"`
	}
	if errParse := c.QueryParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	data, count, errService := service.GetSalesOrders(req.Page, req.Limit, req.TrSalesOrder)
	if errService != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, errService.Error(), nil)
	}

	var resp []model.SalesOrderResponse
	for _, val := range data {
		item := model.ConvertSalesOrderToResponse(val)
		resp = append(resp, item)
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": resp, "total": count})
}

func GetSalesOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, errParse := strconv.Atoi(idStr)
	if errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	item, errService := service.GetSalesOrder(uint(id))
	if errService != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, errService.Error(), nil)
	}

	resp := model.ConvertSalesOrderToResponse(item)
	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", resp)
}

func CreateSalesOrder(c *fiber.Ctx) error {
	var req model.SalesOrderResponse
	if errParse := c.BodyParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	today := time.Now()
	lastSoSequence, _ := service.GetLastSoSequence(today)
	nextSoSequence := lastSoSequence + 1
	soCode := fmt.Sprintf("SO-%s%03d", today.Format("060102"), nextSoSequence)
	req.SoCode = soCode

	var prices []int
	var quantities []int
	for _, val := range req.SalesOrderProducts {
		productPrice, _ := service.GetLastProductPrice(val.ProductID)
		prices = append(prices, productPrice.SellPrice)
		quantities = append(quantities, val.Quantity)
	}
	req.Amount, _ = helper.CountAmount(req.Tax, req.Disc, prices, quantities)

	payload := model.ConvertResponseToSalesOrder(req)
	create, errService := service.CreateSalesOrder(payload)
	if errService != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, errService.Error(), nil)
	}

	if errDecrementStock := decrementStock(req.SalesOrderProducts, req.UserID); errDecrementStock != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, errDecrementStock.Error(), nil)
	}

	find, _ := service.GetSalesOrder(create.ID)
	resp := model.ConvertSalesOrderToResponse(find)
	return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", resp)
}

func decrementStock(products []model.SalesOrderProductResponse, userID uint) error {
	for _, val := range products {
		errAddStock := service.CreateStock(converter.Decrement, userID, val.ProductID, val.Quantity)
		if errAddStock != nil {
			return errAddStock
		}
	}
	return nil
}
