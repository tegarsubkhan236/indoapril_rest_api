package entities

import "errors"

// CrPermission act as table and request body
type CrPermission struct {
	ID       uint           `json:"id" gorm:"primary_key"`
	ParentID *uint          `json:"parent_id"`
	Name     string         `json:"name" gorm:"size:64;not null"`
	Children []CrPermission `json:"-" gorm:"foreignkey:ParentID"`
}

// CrPermissionResp act as response body
type CrPermissionResp struct {
	ID       uint               `json:"id"`
	ParentID *uint              `json:"parent_id,omitempty"`
	Name     string             `json:"name"`
	Children []CrPermissionResp `json:"children,omitempty"`
}

func (r CrPermission) ValidateInput() error {
	if r.Name == "" {
		return errors.New("please specify name")
	}
	return nil
}

func (r CrPermission) ToResponse() CrPermissionResp {
	permission := CrPermissionResp{
		ID:       r.ID,
		ParentID: r.ParentID,
		Name:     r.Name,
	}

	for _, child := range r.Children {
		permission.Children = append(permission.Children, child.ToResponse())
	}

	return permission
}
