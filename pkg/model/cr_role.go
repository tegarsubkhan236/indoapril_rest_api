package model

type CrRole struct {
	ID          uint           `gorm:"primary_key"`
	ParentID    *uint          `gorm:"index"`
	Name        string         `gorm:"size:64;not null"`
	Children    []CrRole       `gorm:"foreignkey:ParentID"`
	Users       []CrUser       `gorm:"foreignKey:RoleID"`
	Permissions []CrPermission `gorm:"many2many:cr_role_permissions"`
}

type CrRoleRequest struct {
	Name        string `json:"name"`
	ParentID    *uint  `json:"parent_id"`
	Users       []uint `json:"users"`
	Permissions []uint `json:"permissions"`
}

type CrRoleResponse struct {
	ID          uint                   `json:"id"`
	ParentID    *uint                  `json:"parent_id,omitempty"`
	Name        string                 `json:"name"`
	Permissions []CrPermissionResponse `json:"permissions,omitempty"`
	Users       []CrUserResponse       `json:"users,omitempty"`
	Children    []CrRoleResponse       `json:"children,omitempty"`
}

func (req CrRoleRequest) ValidateInput() []string {
	var errValidate []string
	if req.Name == "" {
		errValidate = append(errValidate, "name is required")
	}
	if len(req.Permissions) == 0 {
		errValidate = append(errValidate, "role has no permission")
	}
	return errValidate
}

func (req CrRoleRequest) ToModel() CrRole {
	role := CrRole{
		ParentID: req.ParentID,
		Name:     req.Name,
	}

	for _, permissionID := range req.Permissions {
		role.Permissions = append(role.Permissions, CrPermission{
			ID: permissionID,
		})
	}

	for _, userID := range req.Users {
		role.Users = append(role.Users, CrUser{
			ID: userID,
		})
	}

	return role
}

func (model CrRole) ToResponse() CrRoleResponse {
	role := CrRoleResponse{
		ID:       model.ID,
		ParentID: model.ParentID,
		Name:     model.Name,
	}

	for _, permission := range model.Permissions {
		role.Permissions = append(role.Permissions, CrPermissionResponse{
			ID:   permission.ID,
			Name: permission.Name,
		})
	}

	for _, user := range model.Users {
		role.Users = append(role.Users, CrUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Sex:      user.Sex,
			Phone:    user.Phone,
			Status:   user.Status,
			Avatar:   user.Avatar,
		})
	}

	for _, child := range model.Children {
		role.Children = append(role.Children, child.ToResponse())
	}

	return role
}
