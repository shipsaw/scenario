package models

import "github.com/jinzhu/gorm"

type Asset struct {
	gorm.Model
	CreatorID uint   `gorm:"not_null; index"`
	Name      string `gorm:"not_null"`
	URL       string
}
