package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

////////////////////////////////////// Public Interfaces //////////////////////////////////////

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

////////////////// Implementation of interfaces ////////////////////////////////

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGorm{} // Interface implementation check

///////////////////////////////////// galleryService ///////////////////////////////////////////

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

/////////////////////////////////// userValidator functions //////////////////////////////////
////////////////////////////////////////////// userValidator ////////////////////////////////////
/////////////////////////////////// userGorm //////////////////////////////////////

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
