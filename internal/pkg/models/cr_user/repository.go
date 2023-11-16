package cr_user

import (
	"errors"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(payload *entities.CrUser) (*entities.CrUser, error)
	ReadUser(page, limit int) (*[]entities.CrUser, int64, error)
	ReadUserByID(ID uint) (*entities.CrUser, error)
	ReadUserByEmail(email string) (*entities.CrUser, error)
	ReadUserByUsername(username string) (*entities.CrUser, error)
	UpdateUser(item *entities.CrUser, payload *entities.CrUser) (*entities.CrUser, error)
	DeleteUser(ID []uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateUser(payload *entities.CrUser) (*entities.CrUser, error) {
	if err := r.DB.Create(&payload).Error; err != nil {
		return nil, err
	}

	return payload, nil
}

func (r repository) ReadUser(page, limit int) (*[]entities.CrUser, int64, error) {
	var data []entities.CrUser
	var count int64
	var offset = (page - 1) * limit

	r.DB = r.DB.Model(&data)
	r.DB = r.DB.Preload("Team")
	r.DB = r.DB.Preload("Role")
	r.DB = r.DB.Preload("Role.Permissions")

	if err := r.DB.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadUserByID(ID uint) (*entities.CrUser, error) {
	var item entities.CrUser

	r.DB = r.DB.Preload("Team")
	r.DB = r.DB.Preload("Role")
	r.DB = r.DB.Preload("Role.Permissions")

	if err := r.DB.First(&item, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r repository) ReadUserByEmail(email string) (*entities.CrUser, error) {
	var item entities.CrUser

	r.DB = r.DB.Preload("Role")
	r.DB = r.DB.Preload("Role.Permissions")

	if err := r.DB.Where("email = ?", email).First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r repository) ReadUserByUsername(username string) (*entities.CrUser, error) {
	var item entities.CrUser

	if err := r.DB.Where("username = ?", username).First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r repository) UpdateUser(item *entities.CrUser, payload *entities.CrUser) (*entities.CrUser, error) {
	if err := r.DB.Model(&item).Updates(payload).Error; err != nil {
		return nil, err
	}
	return payload, nil
}

func (r repository) DeleteUser(ID []uint) error {
	var data []entities.CrUser
	var tx = r.DB.Begin()

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if err := r.DB.Unscoped().Delete(&data, ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
