package cr_permission

import (
	"example/internal/pkg/entities"
)

type Service interface {
	Insert(payload *entities.CrPermission) error
	FetchAll(page, limit int) (*[]entities.CrPermissionResp, int64, error)
	FetchDetail(ID uint) (*entities.CrPermissionResp, error)
	Update(ID uint, payload *entities.CrPermission) error
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

func (s service) Insert(payload *entities.CrPermission) error {
	return s.repository.Create(payload)
}

func (s service) FetchAll(page, limit int) (*[]entities.CrPermissionResp, int64, error) {
	permissions, count, err := s.repository.ReadAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var results []entities.CrPermissionResp
	for _, item := range *permissions {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetail(ID uint) (*entities.CrPermissionResp, error) {
	payload, err := s.repository.ReadByID(ID)
	if err != nil {
		return nil, err
	}
	response := payload.ToResponse()

	return &response, nil
}

func (s service) Update(ID uint, payload *entities.CrPermission) error {
	item, err := s.repository.ReadByID(ID)
	if err != nil {
		return err
	}

	return s.repository.Update(item, payload)
}

func (s service) Delete(ID uint) error {
	return s.repository.Delete(ID)
}
