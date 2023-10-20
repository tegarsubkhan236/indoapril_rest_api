package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"github.com/gofiber/fiber/v2"
)

func GetPermissions(c *fiber.Ctx) error {
	var query struct {
		Page  int `query:"page" default:"1"`
		Limit int `query:"limit" default:"10"`
	}
	if err := c.QueryParser(&query); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	data, count, err := service.GetAllPermission(query.Page, query.Limit)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	var resp []model.CrPermissionResponse
	for _, item := range data {
		resp = append(resp, item.ToResponse())
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": resp, "total": count})
}

func GetPermission(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	item, err := service.GetPermissionById(id)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", item.ToResponse())
}

func CreatePermission(c *fiber.Ctx) error {
	var req model.CrPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err.Error())
	}

	if err := req.ValidateInput(); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err)
	}

	result, err := service.CreatePermission(req.ToModel())
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusCreated, "Success create data", result.ToResponse())
}

func UpdatePermission(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	var req model.CrPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := req.ValidateInput(); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err)
	}

	result, err := service.UpdatePermission(id, req.ToModel())
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update data", result.ToResponse())
}

func DeletePermission(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	if err := service.DestroyPermission(id); err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success delete data", nil)
}
