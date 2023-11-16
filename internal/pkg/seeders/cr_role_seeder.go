package seeders

import (
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

func CrRoleSeeder(db *gorm.DB) {
	roles := []entities.CrRole{
		{
			ParentID: nil,
			Name:     "SUPER ADMIN",
		},
	}

	if err := db.Create(&roles).Error; err != nil {
		panic("failed to seed CrRole data")
	}

	superAdmin := roles[0]
	permissions := []entities.CrPermission{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
		{
			ID: 3,
		},
		{
			ID: 4,
		},
		{
			ID: 5,
		},
	}

	for _, permission := range permissions {
		err := db.Model(&superAdmin).Association("Permissions").Append(&permission)
		if err != nil {
			panic("failed to associate permission with role")
		}
	}
}
