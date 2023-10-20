package service

import (
	"example/pkg"
	"example/pkg/model"
)

func GetAllPermission(page, limit int) ([]model.CrPermission, int64, error) {
	var db = pkg.DB
	var count int64
	var data []model.CrPermission
	var offset = (page - 1) * limit

	db = db.Model(&data)
	db = db.Preload("Children").Preload("Children.Children")

	if err := db.Count(&count).Where("parent_id IS NULL").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func GetPermissionById(id uint) (*model.CrPermission, error) {
	var db = pkg.DB
	var item model.CrPermission

	db = db.Preload("Children").Preload("Children.Children")

	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func CreatePermission(payload model.CrPermission) (*model.CrPermission, error) {
	var db = pkg.DB
	if err := db.Create(&payload).Error; err != nil {
		return nil, err
	}
	return &payload, nil
}

func UpdatePermission(id uint, payload model.CrPermission) (*model.CrPermission, error) {
	var db = pkg.DB

	item, err := GetPermissionById(id)
	if err != nil {
		return nil, err
	}

	if err := db.Model(&item).Updates(payload).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func DestroyPermission(id uint) error {
	var db = pkg.DB
	result := db.Unscoped().Delete(&model.CrPermission{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
