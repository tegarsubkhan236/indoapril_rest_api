package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"github.com/gofiber/fiber/v2"
)

func GetRoles(c *fiber.Ctx) error {
	var query struct {
		Page   int          `query:"page" default:"1"`
		Limit  int          `query:"limit" default:"10"`
		Filter model.CrRole `query:"filter"`
	}
	if err := c.QueryParser(&query); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	data, count, err := service.GetAllRole(query.Page, query.Limit, query.Filter)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	var resp []model.CrRoleResponse
	for _, item := range data {
		resp = append(resp, item.ToResponse())
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": resp, "total": count})
}

func GetRole(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	item, err := service.GetRoleById(id)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", item.ToResponse())
}

func CreateRole(c *fiber.Ctx) error {
	var req model.CrRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := req.ValidateInput(); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "validate error", err)
	}

	if err := service.CreateRole(req.ToModel()); err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success create data", nil)
}

func UpdateRole(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	var req model.CrRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := req.ValidateInput(); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "validate error", err)
	}

	if err := service.UpdateRole(id, req.ToModel()); err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update data", nil)
}

func DeleteRole(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	if err := service.DestroyRole(id); err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusNoContent, "success delete data", nil)
}
