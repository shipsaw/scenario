package models

type ImageService interface {
	Create() error
	// ByGalleryID(galleryID uint) []string

	func NewImageService(db *gorm.DB) ImageService {
		return &imageService{
			ImageDB: &imageValidator{&imageGorm{db}},
		}
	}
}
