package model

import (
	"gorm.io/gorm"
	"time"
)

type CrTeam struct {
	ID        uint `gorm:"primary_key"`
	ParentID  *uint
	Name      string         `gorm:"size:64;not null"`
	Children  []CrTeam       `gorm:"foreignkey:ParentID"`
	Users     []CrUser       `gorm:"foreignKey:TeamID"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CrTeamRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name"`
	UserIDs  []uint `json:"user_ids"`
}

type CrTeamResponse struct {
	ID       uint             `json:"id"`
	ParentID *uint            `json:"parent_id"`
	Name     string           `json:"name"`
	Children []CrTeamResponse `json:"children,omitempty"`
	Users    []CrUserResponse `json:"users,omitempty"`
}

func (req CrTeamRequest) ValidateInput() []string {
	var errValidate []string
	if req.Name == "" {
		errValidate = append(errValidate, "name is required")
	}
	if len(req.UserIDs) == 0 {
		errValidate = append(errValidate, "team has no member")
	}
	return errValidate
}

func (req CrTeamRequest) ToModel() CrTeam {
	item := CrTeam{
		ParentID: req.ParentID,
		Name:     req.Name,
	}
	for _, x := range req.UserIDs {
		user := CrUser{
			ID: x,
		}
		item.Users = append(item.Users, user)
	}
	return item
}
