package model

import "gorm.io/gorm"

type favorite struct {
	ProductID uint64
	UserID    uint64
	User      User    `gorm:"foreignkey:UserID"`
	Product   Product `gorm:"foreignkey:ProductID"`
	gorm.Model
}
