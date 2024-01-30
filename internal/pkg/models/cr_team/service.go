package cr_team

import (
	"example/internal/pkg/entities"
)

type Service interface {
	Insert(team *entities.CrTeam) error
	FetchAll(page, limit int) (*[]entities.CrTeam, int64, error)
	FetchDetail(ID uint) (*entities.CrTeam, error)
	Update(ID uint, team *entities.CrTeam) error
	Delete(ID uint) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) Insert(team *entities.CrTeam) error {
	//TODO implement me
	panic("implement me")
}

func (s service) FetchAll(page, limit int) (*[]entities.CrTeam, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) FetchDetail(ID uint) (*entities.CrTeam, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Update(ID uint, team *entities.CrTeam) error {
	//TODO implement me
	panic("implement me")
}

func (s service) Delete(ID uint) error {
	//TODO implement me
	panic("implement me")
}
