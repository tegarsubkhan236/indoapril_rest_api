package model

import (
	"gorm.io/gorm"
	"time"
)

type CrUser struct {
	ID        uint           `gorm:"primary_key"`
	Username  string         `gorm:"type:varchar(32);index,unique;not null"`
	Email     string         `gorm:"type:varchar(32);index,unique;not null"`
	Sex       int            `gorm:"type:tinyint"`
	Phone     string         `gorm:"type:varchar(32)"`
	Status    int            `gorm:"type:tinyint"`
	TeamID    uint           `gorm:"index"`
	Team      CrTeam         `gorm:"foreignKey:TeamID"`
	RoleID    uint           `gorm:"index"`
	Role      CrRole         `gorm:"foreignKey:RoleID"`
	Avatar    string         `gorm:"type:varchar(255)"`
	Password  string         `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CrUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Sex      int    `json:"sex"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
	TeamID   uint   `json:"team_id"`
	RoleID   uint   `json:"role_id"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
}

type CrUserResponse struct {
	ID       uint           `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Sex      int            `json:"sex"`
	Phone    string         `json:"phone"`
	Status   int            `json:"status"`
	Avatar   string         `json:"avatar"`
	Team     CrTeamResponse `json:"team,omitempty"`
	Role     CrRoleResponse `json:"role,omitempty"`
}

func (req CrUserRequest) ValidateInput() []string {
	var errValidate []string
	if req.Username == "" {
		errValidate = append(errValidate, "Username is required")
	}
	if req.Email == "" {
		errValidate = append(errValidate, "Email is required")
	}
	if req.Sex != 0 && req.Sex != 1 {
		errValidate = append(errValidate, "Sex should be either 0 or 1")
	}
	if req.Phone == "" {
		errValidate = append(errValidate, "Phone is required")
	}
	if req.Status != 0 && req.Status != 1 {
		errValidate = append(errValidate, "Status should be either 0 or 1")
	}
	if req.Avatar == "" {
		errValidate = append(errValidate, "Avatar is required")
	}
	if req.Password == "" {
		errValidate = append(errValidate, "Password is required")
	}
	if req.RoleID == 0 {
		errValidate = append(errValidate, "RoleID is required")
	}
	if req.TeamID == 0 {
		errValidate = append(errValidate, "TeamID is required")
	}
	return errValidate
}

func (req CrUserRequest) ToModel() CrUser {
	return CrUser{
		Username: req.Username,
		Email:    req.Email,
		Sex:      req.Sex,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Password: req.Password,
		RoleID:   req.RoleID,
		TeamID:   req.TeamID,
		Status:   req.Status,
	}
}

func (model CrUser) ToResponse() CrUserResponse {
	item := CrUserResponse{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Sex:      model.Sex,
		Phone:    model.Phone,
		Avatar:   model.Avatar,
		Status:   model.Status,
		Team: CrTeamResponse{
			ID:   model.Team.ID,
			Name: model.Team.Name,
		},
		Role: CrRoleResponse{
			ID:   model.Role.ID,
			Name: model.Role.Name,
		},
	}
	for _, permission := range model.Role.Permissions {
		item.Role.Permissions = append(item.Role.Permissions, CrPermissionResponse{
			ID:       permission.ID,
			ParentID: permission.ParentID,
			Name:     permission.Name,
		})
	}
	return item
}
