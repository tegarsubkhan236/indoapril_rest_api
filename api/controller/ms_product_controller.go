package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetProducts(c *fiber.Ctx) error {
	var req struct {
		Page            int    `query:"page" default:"1"`
		Limit           int    `query:"limit" default:"10"`
		SupplierID      int    `query:"supplier_id"`
		ProductCategory []int  `query:"product_category"`
		SearchText      string `query:"search_text"`
	}
	if err := c.QueryParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	data, count, err := service.GetAllProduct(req.Page, req.Limit, req.SupplierID, req.ProductCategory, req.SearchText)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	var resp []model.ProductResponse
	for _, item := range data {
		resp = append(resp, item.ToResponse())
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": resp, "total": count})
}

func GetProduct(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	var basedOn = c.Params("based_on") // Price or Stock

	item, err := service.GetProductById(id, basedOn)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", item.ToResponse())
}

func CreateProduct(c *fiber.Ctx) error {
	var req []model.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	var products []model.MsProduct
	var validateCategories []string
	for _, item := range req {
		if len(item.ProductCategoriesID) == 0 {
			validateCategory := fmt.Sprintf("Product %s has no category", item.Name)
			validateCategories = append(validateCategories, validateCategory)
		}
		products = append(products, item.ToModel())
	}
	if len(validateCategories) > 0 {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "validate category", validateCategories)
	}

	createResult, err := service.CreateProduct(products)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	if len(createResult) == 1 {
		item, _ := service.GetProductById(createResult[0].ID, "")
		return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", item.ToResponse())
	}

	return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", len(createResult))
}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}
	var payload model.MsProduct

	if err := c.BodyParser(&payload); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	item, err := service.GetProductById(id, "")
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	result, err := service.UpdateProduct(*item, payload)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update data", fiber.Map{"user": result})
}

func DeleteProduct(c *fiber.Ctx) error {
	var req struct {
		ID []int `query:"id"`
	}
	if err := c.QueryParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err := service.DestroyProduct(req.ID)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusNoContent, "success delete data", nil)
}
