package cr_role

import (
	"example/internal/pkg/entities"
)

type Service interface {
	Insert(role *entities.CrRole) error
	FetchAll(page, limit int) (*[]entities.CrRoleResp, int64, error)
	FetchDetail(ID uint) (*entities.CrRoleResp, error)
	Update(ID uint, role *entities.CrRole) error
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

func (s service) Insert(role *entities.CrRole) error {
	return s.repository.Create(role)
}

func (s service) FetchAll(page, limit int) (*[]entities.CrRoleResp, int64, error) {
	roles, count, err := s.repository.ReadAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var results []entities.CrRoleResp
	for _, item := range *roles {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetail(ID uint) (*entities.CrRoleResp, error) {
	role, err := s.repository.ReadByID(ID)
	if err != nil {
		return nil, err
	}
	response := role.ToResponse()

	return &response, nil
}

func (s service) Update(ID uint, role *entities.CrRole) error {
	oldRole, err := s.repository.ReadByID(ID)
	if err != nil {
		return err
	}

	return s.repository.Update(oldRole, role)
}

func (s service) Delete(ID uint) error {
	return s.repository.Delete(ID)
}
