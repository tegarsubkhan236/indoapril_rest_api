package model

type CrPermission struct {
	ID       uint `gorm:"primary_key"`
	ParentID *uint
	Name     string         `gorm:"size:64;not null"`
	Children []CrPermission `gorm:"foreignkey:ParentID"`
}

type CrPermissionRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name"`
}

type CrPermissionResponse struct {
	ID       uint                   `json:"id"`
	ParentID *uint                  `json:"parent_id,omitempty"`
	Name     string                 `json:"name"`
	Children []CrPermissionResponse `json:"children,omitempty"`
}

func (req CrPermissionRequest) ValidateInput() []string {
	var errValidate []string
	if req.Name == "" {
		errValidate = append(errValidate, "name is required")
	}
	return errValidate
}

func (req CrPermissionRequest) ToModel() CrPermission {
	return CrPermission{
		ParentID: req.ParentID,
		Name:     req.Name,
	}
}

func (model CrPermission) ToResponse() CrPermissionResponse {
	permission := CrPermissionResponse{
		ID:       model.ID,
		ParentID: model.ParentID,
		Name:     model.Name,
	}

	for _, child := range model.Children {
		permission.Children = append(permission.Children, child.ToResponse())
	}

	return permission
}
