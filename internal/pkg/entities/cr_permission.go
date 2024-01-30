package entities

import "errors"

// CrPermission act as table, request body, response body
type CrPermission struct {
	ID       uint           `json:"id" gorm:"primary_key"`
	ParentID *uint          `json:"parent_id"`
	Name     string         `json:"name" gorm:"size:64;not null"`
	Children []CrPermission `json:"children" gorm:"foreignkey:ParentID"`
}

func (r CrPermission) ValidateInput() error {
	if r.Name == "" {
		return errors.New("please specify name")
	}
	return nil
}

func (r CrPermission) ToResponse() CrPermission {
	item := CrPermission{
		ID:       r.ID,
		ParentID: r.ParentID,
		Name:     r.Name,
	}

	for _, child := range r.Children {
		item.Children = append(item.Children, child.ToResponse())
	}

	return item
}
