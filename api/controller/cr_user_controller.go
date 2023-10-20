package controller

import (
	"example/api/service"
	"example/api/tool/helper"
	"example/pkg/model"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var query struct {
		Page          int                 `query:"page" default:"1"`
		Limit         int                 `query:"limit" default:"10"`
		CrUserRequest model.CrUserRequest `query:"user"`
	}
	if err := c.QueryParser(&query); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err.Error())
	}

	data, count, err := service.GetAllUser(query.Page, query.Limit, query.CrUserRequest)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	var resp []model.CrUserResponse
	for _, item := range data {
		resp = append(resp, item.ToResponse())
	}
	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", fiber.Map{"results": resp, "total": count})
}

func GetUser(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	user, err := service.GetUserById(id)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success find data", user.ToResponse())
}

func CreateUser(c *fiber.Ctx) error {
	var req model.CrUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err.Error())
	}

	hash, err := service.HashPassword(req.Password)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "Couldn't hash password", err.Error())
	}

	req.Password = hash
	result, err := service.CreateUser(req.ToModel())
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "Couldn't create user", err.Error())
	}

	return helper.ResponseHandler(c, fiber.StatusCreated, "Created user", result)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	var req model.CrUserRequest
	if errParsePayload := c.BodyParser(&req); errParsePayload != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", errParsePayload.Error())
	}

	user, errGet := service.GetUserById(id)
	if errGet != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, "No data find with ID", errGet.Error())
	}

	avatarOldPath := user.Avatar
	avatarNewPath, errUpload := helper.SingleUpload(c, "avatar", avatarOldPath)
	if errUpload != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "Failed to handle avatar upload", err)
	}
	if avatarNewPath != "" {
		req.Avatar = avatarNewPath
	}

	result, errUpdate := service.UpdateUser(*user, req.ToModel())
	if errUpdate != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, "Failed to update data", errUpdate.Error())
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update data", result)
}

func UpdateUserPassword(c *fiber.Ctx) error {
	id, err := helper.GetIDParam(c)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, err.Error(), nil)
	}

	var req struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Review your input", err.Error())
	}
	if req.NewPassword != req.ConfirmPassword {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Password confirmation does not match", nil)
	}

	user, err := service.GetUserById(id)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, "User not found", err.Error())
	}

	if !service.CheckPasswordHash(req.OldPassword, user.Password) {
		return helper.ResponseHandler(c, fiber.StatusNotFound, "Old password is incorrect", nil)
	}

	hashedNewPass, _ := service.HashPassword(req.NewPassword)
	var newUserPayload model.CrUser
	newUserPayload.Password = hashedNewPass
	result, err := service.UpdateUser(*user, newUserPayload)
	if err != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, "Failed to update data", err.Error())
	}

	return helper.ResponseHandler(c, fiber.StatusOK, "Success update password", result)
}

func DeleteUser(c *fiber.Ctx) error {
	var batchQuery struct {
		ID []int `query:"id"`
	}
	if err := c.QueryParser(&batchQuery); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "Failed to parse request query", err.Error())
	}

	if err := service.DestroyUser(batchQuery.ID); err != nil {
		return helper.ResponseHandler(c, fiber.StatusBadGateway, "error", err.Error())
	}

	return helper.ResponseHandler(c, fiber.StatusNoContent, "", nil)
}
