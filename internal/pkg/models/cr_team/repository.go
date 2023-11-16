package cr_team

import (
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	Create(payload *entities.CrTeam) (*entities.CrTeam, error)
	ReadAll(ID uint, filer *entities.CrTeam, page, limit int) (*[]entities.CrTeam, int64, error)
	Update(ID uint, payload *entities.CrTeam) (*entities.CrTeam, error)
	Delete(ID uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) Create(payload *entities.CrTeam) (*entities.CrTeam, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) ReadAll(ID uint, filer *entities.CrTeam, page, limit int) (*[]entities.CrTeam, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(ID uint, payload *entities.CrTeam) (*entities.CrTeam, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ID uint) error {
	//TODO implement me
	panic("implement me")
}
