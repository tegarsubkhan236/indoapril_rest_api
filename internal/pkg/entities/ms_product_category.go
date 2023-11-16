package entities

import (
	"gorm.io/gorm"
	"time"
)

type MsProductCategory struct {
	ID        uint                `gorm:"primary_key" json:"id"`
	ParentID  *uint               `json:"parent_id"`
	Name      string              `gorm:"size:32;not null" json:"name"`
	Children  []MsProductCategory `gorm:"foreignkey:ParentID" json:"children,omitempty"`
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"-"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
}

type MsProductCategoryResp struct {
	ID       uint                    `json:"id"`
	ParentID *uint                   `json:"parent_id"`
	Name     string                  `json:"name"`
	Children []MsProductCategoryResp `json:"children,omitempty"`
}

func (r MsProductCategory) ToResponse() MsProductCategoryResp {
	productCategory := MsProductCategoryResp{
		ID:       r.ID,
		ParentID: r.ParentID,
		Name:     r.Name,
	}

	for _, child := range r.Children {
		productCategory.Children = append(productCategory.Children, child.ToResponse())
	}

	return productCategory
}
