package seeder

import (
	"example/pkg/model"
	"gorm.io/gorm"
)

func CrRoleSeeder(db *gorm.DB) {
	roles := []model.CrRole{
		{
			ParentID: nil,
			Name:     "SUPER ADMIN",
		},
	}

	result := db.Create(&roles)
	if result.Error != nil {
		panic("failed to seed CoreRole data")
	}
}
