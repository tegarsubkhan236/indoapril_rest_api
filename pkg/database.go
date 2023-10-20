package pkg

import (
	"example/pkg/model"
	"example/pkg/seeder"
	"fmt"
	//"gorm.io/driver/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", GetEnv("DB_USER"), GetEnv("DB_PASSWORD"), GetEnv("DB_HOST"), GetEnv("DB_PORT"), GetEnv("DB_NAME"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("can't connect to database")
	}
	fmt.Println("Connection success")

	RunMigration(DB)
	RunSeeder()
}

func RunMigration(DB *gorm.DB) {
	err := DB.AutoMigrate(
		&model.CrSetting{},
		&model.CrPermission{},
		&model.CrRole{},
		&model.CrUser{},
		&model.CrTeam{},
		&model.MsSupplier{},
		&model.MsProductCategory{},
		&model.MsProduct{},
		&model.MsProductPrice{},
		&model.MsStock{},
		&model.TrPurchaseOrder{},
		&model.TrPurchaseOrderProduct{},
		&model.TrReceivingOrder{},
		&model.TrReceivingOrderProduct{},
		&model.TrBackOrder{},
		&model.TrBackOrderProduct{},
		&model.TrSalesOrder{},
		&model.TrSalesOrderProduct{},
		&model.TrReturnOrder{},
		&model.TrReturnOrderProduct{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}

func RunSeeder() {
	var coreSetting model.CrSetting

	err := DB.First(&coreSetting).Error
	if err == gorm.ErrRecordNotFound {
		result := DB.Create(&coreSetting)
		if result.Error != nil {
			panic("failed to seed CoreSetting data")
		}
	}

	if coreSetting.IsSeeded {
		return
	}

	seeder.CrRoleSeeder(DB)
	seeder.CrTeamSeeder(DB)
	seeder.CrUserSeeder(DB)

	result := DB.Model(&coreSetting).Where("id = ?", 1).Update("is_seeded", true)
	if result.Error != nil {
		panic("failed to update seeder status")
	}
}
