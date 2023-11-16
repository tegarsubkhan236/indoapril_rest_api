package seeders

import (
	"example/internal/api/util/constant"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

func CrPermissionSeeder(db *gorm.DB) {
	permissions := []entities.CrPermission{
		{
			Name: constant.MANAGE_PERMISSION,
			Children: []entities.CrPermission{
				{
					Name: constant.READ_PERMISSION,
				},
				{
					Name: constant.CREATE_PERMISSION,
				},
				{
					Name: constant.CREATE_PERMISSION,
				},
				{
					Name: constant.DELETE_PERMISSION,
				},
			},
		},
		{
			Name: constant.MANAGE_ROLE,
			Children: []entities.CrPermission{
				{
					Name: constant.READ_ROLE,
				},
				{
					Name: constant.CREATE_ROLE,
				},
				{
					Name: constant.UPDATE_ROLE,
				},
				{
					Name: constant.DELETE_ROLE,
				},
			},
		},
		{
			Name: constant.MANAGE_USER,
			Children: []entities.CrPermission{
				{
					Name: constant.READ_USER,
				},
				{
					Name: constant.CREATE_USER,
				},
				{
					Name: constant.UPDATE_USER,
				},
				{
					Name: constant.DELETE_USER,
				},
			},
		},
	}

	result := db.Create(&permissions)
	if result.Error != nil {
		panic("failed to seed CrPermission data")
	}
}
