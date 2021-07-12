package models

import "github.com/jinzhu/gorm"

type Services struct {
	db      *gorm.DB
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
		db:      db,
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
	}, nil
}

// Drops and rebuilds all tables
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// Wrapper around Gorm automigrate to allow us to be db-type agnostic
func (s *Services) AutoMigrate() error {
	err := s.db.AutoMigrate(&User{}, &Gallery{}).Error
	return err
}

func (s *Services) Close() error {
	return s.db.Close()
}
