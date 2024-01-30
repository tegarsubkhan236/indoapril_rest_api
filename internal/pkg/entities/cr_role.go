package entities

// CrRole act as table and response body
type CrRole struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	ParentID    *uint          `json:"parent_id" gorm:"index"`
	Name        string         `json:"name" gorm:"size:64;not null"`
	Children    []CrRole       `json:"children" gorm:"foreignkey:ParentID"`
	Users       []CrUser       `json:"users" gorm:"foreignKey:RoleID"`
	Permissions []CrPermission `json:"permissions" gorm:"many2many:cr_role_permissions"`
}

// CrRoleReq act as request body
type CrRoleReq struct {
	ID            uint   `json:"id"`
	ParentID      *uint  `json:"parent_id"`
	Name          string `json:"name"`
	PermissionIds []int  `json:"permissions_ids"`
}

func (r CrRole) ToResponse() CrRole {
	role := CrRole{
		ID:       r.ID,
		ParentID: r.ParentID,
		Name:     r.Name,
	}

	for _, permission := range r.Permissions {
		role.Permissions = append(role.Permissions, CrPermission{
			ID:   permission.ID,
			Name: permission.Name,
		})
	}

	for _, user := range r.Users {
		role.Users = append(role.Users, CrUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Sex:      user.Sex,
			Phone:    user.Phone,
			Status:   user.Status,
			Avatar:   user.Avatar,
		})
	}

	for _, child := range r.Children {
		role.Children = append(role.Children, child.ToResponse())
	}

	return role
}
