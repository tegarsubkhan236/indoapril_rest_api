package service

import (
	"example/pkg"
	"example/pkg/model"
)

func GetAllRole(page, limit int, filter model.CrRole) ([]model.CrRole, int64, error) {
	var db = pkg.DB
	var data []model.CrRole
	var count int64
	var offset = (page - 1) * limit

	db = db.Model(&data)
	db = db.Preload("Permissions")
	db = db.Preload("Children").Preload("Children.Children")

	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := db.Count(&count).Where("parent_id IS NULL").Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func GetRoleById(id uint) (*model.CrRole, error) {
	var db = pkg.DB
	var item model.CrRole

	db = db.Model(&item)
	db = db.Preload("Permissions")
	db = db.Preload("Children").Preload("Children.Children")
	db = db.Preload("Children.Permissions").Preload("Children.Children.Permissions")

	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func CreateRole(payload model.CrRole) error {
	var db = pkg.DB

	if err := db.Create(&payload).Error; err != nil {
		return err
	}

	return nil
}

func UpdateRole(id uint, payload model.CrRole) error {
	var db = pkg.DB
	var tx = db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	item, err := GetRoleById(id)
	if err != nil {
		return err
	}

	var permissions []model.CrPermission
	for _, permission := range payload.Permissions {
		var p model.CrPermission
		if err := db.First(&p, permission.ID).Error; err != nil {
			return err
		}
		permissions = append(permissions, p)
	}

	if err := db.Model(&item).Updates(payload).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := db.Model(&item).Association("Permissions").Replace(permissions); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DestroyRole(id uint) error {
	var db = pkg.DB
	var role model.CrRole
	var tx = db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := db.First(&role, id).Error; err != nil {
		return err
	}

	if err := db.Model(&role).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	if err := db.Unscoped().Delete(&role).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
