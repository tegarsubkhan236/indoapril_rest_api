package seeders

import (
	"example/internal/api/types/permissions"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

func CrPermissionSeeder(db *gorm.DB) {
	crPermissions := []entities.CrPermission{
		{
			Name: permissions.MANAGE_PERMISSION,
			Children: []entities.CrPermission{
				{
					Name: permissions.READ_PERMISSION,
				},
				{
					Name: permissions.CREATE_PERMISSION,
				},
				{
					Name: permissions.UPDATE_PERMISSION,
				},
				{
					Name: permissions.DELETE_PERMISSION,
				},
			},
		},
		{
			Name: permissions.MANAGE_ROLE,
			Children: []entities.CrPermission{
				{
					Name: permissions.READ_ROLE,
				},
				{
					Name: permissions.CREATE_ROLE,
				},
				{
					Name: permissions.UPDATE_ROLE,
				},
				{
					Name: permissions.DELETE_ROLE,
				},
			},
		},
		{
			Name: permissions.MANAGE_USER,
			Children: []entities.CrPermission{
				{
					Name: permissions.READ_USER,
				},
				{
					Name: permissions.CREATE_USER,
				},
				{
					Name: permissions.UPDATE_USER,
				},
				{
					Name: permissions.DELETE_USER,
				},
			},
		},
		{
			Name: permissions.MANAGE_TEAM,
			Children: []entities.CrPermission{
				{
					Name: permissions.READ_TEAM,
				},
				{
					Name: permissions.CREATE_TEAM,
				},
				{
					Name: permissions.UPDATE_TEAM,
				},
				{
					Name: permissions.DELETE_TEAM,
				},
			},
		},
	}

	result := db.Create(&crPermissions)
	if result.Error != nil {
		panic("failed to seed CrPermission data")
	}
}
