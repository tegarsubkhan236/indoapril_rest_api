package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"github.com/gofiber/fiber/v2"
)

func GetProductCategories(c *fiber.Ctx) error {
	var query struct {
		Page  int `query:"page" default:"1"`
		Limit int `query:"limit" default:"10"`
	}
	if err := c.QueryParser(&query); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err.Error())
	}

	data, count, err := service.GetAllProductCategory(query.Page, query.Limit)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": data, "total": count})
}

func GetProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := service.GetProductCategoryById(id)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", item)
}

func CreateProductCategory(c *fiber.Ctx) error {
	var req model.MsProductCategory
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err := service.CreateProductCategory(req)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", nil)
}

func UpdateProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.MsProductCategory
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	item, err := service.GetProductCategoryById(id)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if err = service.UpdateProductCategory(*item, req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update data", nil)
}

func DeleteProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := service.DestroyProductCategory(id); err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success delete data", nil)
}
