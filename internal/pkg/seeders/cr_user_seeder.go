package seeders

import (
	"example/internal/pkg/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CrUserSeeder(db *gorm.DB) {
	initialPassword := "whiplash"
	hash, err := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to seed CoreUser data")
	}

	users := []entities.CrUser{
		{
			Username: "tegarsubkhan",
			Email:    "wwww@gmail.com",
			Password: string(hash),
			Sex:      1,
			Status:   1,
			RoleID:   1,
			TeamID:   1,
		},
	}

	result := db.Create(&users)
	if result.Error != nil {
		panic("failed to seed CoreUser data")
	}
}
