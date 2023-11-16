package seeders

import (
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

func CrTeamSeeder(db *gorm.DB) {
	teams := []entities.CrTeam{
		{
			Name: "Central Team",
		},
	}

	result := db.Create(&teams)
	if result.Error != nil {
		panic("failed to seed CoreRole data")
	}
}
