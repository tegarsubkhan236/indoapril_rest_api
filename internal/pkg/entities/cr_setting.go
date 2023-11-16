package entities

type CrSetting struct {
	ID       uint `gorm:"primary_key" json:"id"`
	IsSeeded bool `gorm:"column:is_seeded;"`
}
