package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func GetPurchaseOrders(c *fiber.Ctx) error {
	var req struct {
		Page            int                   `query:"page" default:"1"`
		Limit           int                   `query:"limit" default:"10"`
		TrPurchaseOrder model.TrPurchaseOrder `query:"column"`
	}
	if errParse := c.QueryParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	data, count, errService := service.GetPurchaseOrders(req.Page, req.Limit, req.TrPurchaseOrder)
	if errService != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, errService.Error(), nil)
	}

	var resp []model.PurchaseOrderResponse
	for _, val := range data {
		item := model.ConvertPurchaseOrderToResponse(val)
		resp = append(resp, item)
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": resp, "total": count})
}

func GetPurchaseOrder(c *fiber.Ctx) error {
	var req struct {
		PoID   uint   `query:"po_id"`
		BoID   uint   `query:"bo_id"`
		PoCode string `query:"po_code"`
		BoCode string `query:"bo_code"`
	}
	if errParse := c.QueryParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	switch {
	case req.PoID != 0:
		item, _, errService := service.GetPurchaseOrder(req.PoID, 0, "", "")
		if errService != nil {
			return helper.ResponseHandler(c, fiber.StatusNotFound, errService.Error(), nil)
		}
		return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", model.ConvertPurchaseOrderToResponse(item))

	case req.BoID != 0:
		_, item, errService := service.GetPurchaseOrder(0, req.BoID, "", "")
		if errService != nil {
			return helper.ResponseHandler(c, fiber.StatusNotFound, errService.Error(), nil)
		}
		return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", model.ConvertBackOrderToResponse(item))

	case req.PoCode != "" && req.BoCode == "":
		item, _, errService := service.GetPurchaseOrder(0, 0, req.PoCode, "")
		if errService != nil {
			return helper.ResponseHandler(c, fiber.StatusNotFound, errService.Error(), nil)
		}
		return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", model.ConvertPurchaseOrderToResponse(item))

	case req.PoCode != "" && req.BoCode != "":
		_, item, errService := service.GetPurchaseOrder(0, 0, req.PoCode, req.BoCode)
		if errService != nil {
			return helper.ResponseHandler(c, fiber.StatusNotFound, errService.Error(), nil)
		}
		return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", model.ConvertBackOrderToResponse(item))

	default:
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "", nil)
	}
}

func CreatePurchaseOrder(c *fiber.Ctx) error {
	var req model.PurchaseOrderResponse
	if errParse := c.BodyParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	today := time.Now()
	lastPoSequence, _ := service.GetLastPoSequence(today)
	nextPoSequence := lastPoSequence + 1
	poCode := fmt.Sprintf("PO-%s%03d", today.Format("060102"), nextPoSequence)
	req.PoCode = poCode

	var prices []int
	var quantities []int
	for _, val := range req.PurchaseOrderProducts {
		prices = append(prices, val.Price)
		quantities = append(quantities, val.Quantity)
	}
	req.Amount, _ = helper.CountAmount(req.Tax, req.Disc, prices, quantities)

	payload := model.ConvertResponseToPurchaseOrder(req)
	create, errService := service.CreatePurchaseOrder(payload)
	if errService != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, errService.Error(), nil)
	}

	find, _, _ := service.GetPurchaseOrder(create.ID, 0, "", "")
	resp := model.ConvertPurchaseOrderToResponse(find)
	return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", resp)
}
