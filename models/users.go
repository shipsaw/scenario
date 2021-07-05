package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound = errors.New("models: resource not found")
)

// Contains the specific type of DB that we are working with
type UserService struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Email string `gorm:"not null;unique_index"`
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

func (us *UserService) Close() error {
	return us.db.Close()
}

func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.First(&user, "id=?", id).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Drops and rebuilds user table
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}
