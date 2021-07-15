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
	ByID(id uint) (*Gallery, error)
	Update(gallery *Gallery) error
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

/////////////////////////////////// GalleryValidator functions //////////////////////////////////

type galleryValFunc func(*Gallery) error

func runGalleryValFuncs(gallery *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(gallery *Gallery) error {
	if gallery.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

////////////////////////////////////////////// GalleryValidator ////////////////////////////////////

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

/////////////////////////////////// userGorm //////////////////////////////////////

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id=?", id)
	err := first(db, &gallery)
	return &gallery, err
}
