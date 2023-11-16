package pkg

import (
	"errors"
	"example/internal/pkg/entities"
	"example/internal/pkg/seeders"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entities.CrSetting{},
		&entities.CrPermission{},
		&entities.CrRole{},
		&entities.CrUser{},
		&entities.CrTeam{},
		&entities.MsSupplier{},
		&entities.MsProductCategory{},
		&entities.MsProduct{},
		&entities.MsProductPrice{},
		&entities.MsStock{},
		&entities.TrPurchaseOrder{},
		&entities.TrPurchaseOrderProduct{},
		&entities.TrReceivingOrder{},
		&entities.TrReceivingOrderProduct{},
		&entities.TrBackOrder{},
		&entities.TrBackOrderProduct{},
		&entities.TrSalesOrder{},
		&entities.TrSalesOrderProduct{},
		&entities.TrReturnOrder{},
		&entities.TrReturnOrderProduct{},
	); err != nil {
		return err
	}
	return nil
}

func RunSeeder(db *gorm.DB) error {
	var coreSetting entities.CrSetting

	if err := db.First(&coreSetting).Error; err == gorm.ErrRecordNotFound {
		if err = db.Create(&coreSetting).Error; err != nil {
			return errors.New("failed to seed CoreSetting data")
		}
	}

	if coreSetting.IsSeeded {
		return nil
	}

	seeders.CrPermissionSeeder(db)
	seeders.CrRoleSeeder(db)
	seeders.CrTeamSeeder(db)
	seeders.CrUserSeeder(db)

	if err := db.Model(&coreSetting).Where("id = ?", 1).Update("is_seeded", true).Error; err != nil {
		return errors.New("failed to update seeders status")
	}

	return nil
}
