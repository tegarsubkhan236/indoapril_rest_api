package cr_role

import (
	"example/internal/pkg/entities"
)

type Service interface {
	Insert(payload *entities.CrRoleReq) error
	FetchAll(page, limit int) (*[]entities.CrRole, int64, error)
	FetchDetail(ID uint) (*entities.CrRole, error)
	Update(ID uint, payload *entities.CrRoleReq) error
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

func (s service) Insert(payload *entities.CrRoleReq) error {
	return s.repository.Create(payload)
}

func (s service) FetchAll(page, limit int) (*[]entities.CrRole, int64, error) {
	roles, count, err := s.repository.ReadAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var results []entities.CrRole
	for _, item := range *roles {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetail(ID uint) (*entities.CrRole, error) {
	role, err := s.repository.ReadByID(ID)
	if err != nil {
		return nil, err
	}
	response := role.ToResponse()

	return &response, nil
}

func (s service) Update(ID uint, payload *entities.CrRoleReq) error {
	item, err := s.repository.ReadByID(ID)
	if err != nil {
		return err
	}

	return s.repository.Update(item, payload)
}

func (s service) Delete(ID uint) error {
	return s.repository.Delete(ID)
}
