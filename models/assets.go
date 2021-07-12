package models

import "github.com/jinzhu/gorm"

type Asset struct {
	gorm.Model
	CreatorID    uint   `gorm:"not_null; index"`
	SearchName   string `goprm:"not_null"`
	Name         string `gorm:"not_null"`
	Path         string
	URL          string
	Dependancies []Asset
}
