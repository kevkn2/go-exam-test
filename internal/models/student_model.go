package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Name   string `gorm:"size:255;not null;"`
	School string `gorm:"size:255;not null;"`
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}
