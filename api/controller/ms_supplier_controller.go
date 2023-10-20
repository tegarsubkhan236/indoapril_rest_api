package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetSuppliers(c *fiber.Ctx) error {
	var req struct {
		Page     int              `query:"page" default:"1"`
		Limit    int              `query:"limit" default:"10"`
		IDs      []int            `query:"ids"`
		Supplier model.MsSupplier `query:"supplier"`
	}
	if err := c.QueryParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	data, count, err := service.GetAllSupplier(req.Page, req.Limit, req.IDs, req.Supplier)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": data, "total": count})
}

func GetSupplier(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, _ := strconv.Atoi(idStr)
	item, err := service.GetSupplierById(uint(id))
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", item)
}

func CreateSupplier(c *fiber.Ctx) error {
	var req []model.MsSupplier
	if errParse := c.BodyParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	createResult, errCreate := service.CreateSupplier(req)
	if errCreate != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, errCreate.Error(), nil)
	}

	if len(createResult) == 1 {
		item, _ := service.GetSupplierById(createResult[0].ID)
		return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", item)
	}

	return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", nil)
}

func UpdateSupplier(c *fiber.Ctx) error {
	var req model.MsSupplier
	idStr := c.Params("id")
	id, _ := strconv.Atoi(idStr)
	if errParse := c.BodyParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, errParse.Error(), nil)
	}

	item, err := service.GetSupplierById(uint(id))
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	result, err := service.UpdateSupplier(item, req)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update data", fiber.Map{"item": result})
}

func DeleteSupplier(c *fiber.Ctx) error {
	var req struct {
		ID []int `query:"id"`
	}
	if err := c.QueryParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := service.DestroySupplier(req.ID); err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusNoContent, "success delete data", nil)
}
