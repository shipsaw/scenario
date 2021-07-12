package models

import "github.com/jinzhu/gorm"

type Services struct {
	User    UserService
	Gallery GalleryService
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User: NewUserService(db),
	}, nil
}
