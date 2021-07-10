package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shipsaw/scenario/hash"
	"github.com/shipsaw/scenario/rand"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidID       = errors.New("models: id provided is invalid")
	ErrInvalidPassword = errors.New("models: invalid password provided")
)

const userPasswordPepper string = "TNhYZuUBK0"
const hmacSecretKey string = "secret-hmac-key"

type User struct {
	gorm.Model
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; unique_index"`
}

////////////////////////////////////// Public Interfaces //////////////////////////////////////

// models package API
type UserService interface {
	UserDB
	// Verifies provided email and password are correct
	Authenticate(email, password string) (*User, error)
}

// Required for DB interaction
type UserDB interface {
	// Query methods
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// User altering methods
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close db connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

////////////////// Implementation of interfaces ////////////////////////////////
type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// Interface fulfullment check
var _ UserDB = &userGorm{}

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
}

///////////////////////////////////// userService ///////////////////////////////////////////

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPasswordPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

/////////////////////////////////// userGorm //////////////////////////////////////
func NewUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(false)
	return &userGorm{
		db:   db,
		hmac: hash.NewHMAC(hmacSecretKey),
	}, nil
}

func (ug *userGorm) Close() error {
	return ug.db.Close()
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id=?", id)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := ug.hmac.Hash(token)
	db := ug.db.Where("remember_hash=?", hashedToken)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (ug *userGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + userPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Drops and rebuilds user table
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// Wrapper around Gorm automigrate to allow us to be db-type agnostic
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// Wrapper for gorm's First method to check for our custom errors
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else {
		return err
	}
}
