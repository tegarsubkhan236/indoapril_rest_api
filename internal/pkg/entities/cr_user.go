package entities

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type CrUser struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	Username  string         `json:"username" gorm:"type:varchar(32);index,unique;not null"`
	Email     string         `json:"email" gorm:"type:varchar(32);index,unique;not null"`
	Sex       int            `json:"sex" gorm:"type:tinyint"`
	Phone     string         `json:"phone" gorm:"type:varchar(32)"`
	Status    int            `json:"status" gorm:"type:tinyint"`
	TeamID    uint           `json:"team_id" gorm:"index"`
	Team      CrTeam         `json:"-" gorm:"foreignKey:TeamID"`
	RoleID    uint           `json:"role_id" gorm:"index"`
	Role      CrRole         `json:"-" gorm:"foreignKey:RoleID"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255)"`
	Password  string         `json:"-" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CrUserResp struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Sex      int        `json:"sex"`
	Phone    string     `json:"phone"`
	Status   int        `json:"status"`
	Avatar   string     `json:"avatar"`
	Team     CrTeamResp `json:"team"`
	Role     CrRoleResp `json:"role"`
}

func (r CrUser) ValidateInput() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Sex != 0 && r.Sex != 1 {
		return errors.New("sex should be either 0 or 1")
	}
	if r.Phone == "" {
		return errors.New("phone is required")
	}
	if r.Status != 0 && r.Status != 1 {
		return errors.New("status should be either 0 or 1")
	}
	if r.Avatar == "" {
		return errors.New("avatar is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.RoleID == 0 {
		return errors.New("RoleID is required")
	}
	if r.TeamID == 0 {
		return errors.New("TeamID is required")
	}
	return nil
}

func (r CrUser) ToResponse() CrUserResp {
	item := CrUserResp{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
		Sex:      r.Sex,
		Phone:    r.Phone,
		Avatar:   r.Avatar,
		Status:   r.Status,
		Team: CrTeamResp{
			ID:   r.Team.ID,
			Name: r.Team.Name,
		},
		Role: CrRoleResp{
			ID:   r.Role.ID,
			Name: r.Role.Name,
		},
	}
	for _, permission := range r.Role.Permissions {
		item.Role.Permissions = append(item.Role.Permissions, CrPermissionResp{
			ID:       permission.ID,
			ParentID: permission.ParentID,
			Name:     permission.Name,
		})
	}

	return item
}
