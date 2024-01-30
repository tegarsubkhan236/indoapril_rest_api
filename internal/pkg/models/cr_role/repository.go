package cr_role

import (
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	Create(payload *entities.CrRoleReq) error
	ReadAll(page, limit int) (*[]entities.CrRole, int64, error)
	ReadByID(ID uint) (*entities.CrRole, error)
	Update(item *entities.CrRole, payload *entities.CrRoleReq) error
	Delete(ID uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) Create(payload *entities.CrRoleReq) error {
	var tx = r.DB.Begin()

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var permissions []entities.CrPermission
	for _, permissionID := range payload.PermissionIds {
		var p entities.CrPermission
		if err := r.DB.First(&p, permissionID).Error; err != nil {
			return err
		}
		permissions = append(permissions, p)
	}

	if err := r.DB.Model(&payload).Create(&payload).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.DB.Model(&payload).Association("Permissions").Replace(permissions); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r repository) ReadAll(page, limit int) (*[]entities.CrRole, int64, error) {
	var data []entities.CrRole
	var count int64
	var offset = (page - 1) * limit

	r.DB = r.DB.Model(&data)
	r.DB = r.DB.Preload("Permissions")
	r.DB = r.DB.Preload("Children").Preload("Children.Children")

	if err := r.DB.Count(&count).Where("parent_id IS NULL").Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadByID(ID uint) (*entities.CrRole, error) {
	var item entities.CrRole

	r.DB = r.DB.Model(&item)
	r.DB = r.DB.Preload("Permissions")
	r.DB = r.DB.Preload("Children").Preload("Children.Children")
	r.DB = r.DB.Preload("Children.Permissions").Preload("Children.Children.Permissions")

	if err := r.DB.First(&item, ID).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r repository) Update(item *entities.CrRole, payload *entities.CrRoleReq) error {
	var tx = r.DB.Begin()

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var permissions []entities.CrPermission
	for _, permissionID := range payload.PermissionIds {
		var p entities.CrPermission
		if err := r.DB.First(&p, permissionID).Error; err != nil {
			return err
		}
		permissions = append(permissions, p)
	}

	if err := r.DB.Model(&item).Updates(&payload).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.DB.Model(&item).Association("Permissions").Replace(permissions); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r repository) Delete(ID uint) error {
	var role entities.CrRole
	var tx = r.DB.Begin()

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if err := r.DB.First(&role, ID).Error; err != nil {
		return err
	}

	if err := r.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.DB.Unscoped().Delete(&role).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
