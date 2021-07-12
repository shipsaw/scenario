package models

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shipsaw/scenario/hash"
	"github.com/shipsaw/scenario/rand"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound          = errors.New("models: resource not found")
	ErrIDInvalid         = errors.New("models: id provided is invalid")
	ErrPasswordIncorrect = errors.New("models: incorrect password provided")
	ErrEmailInvalid      = errors.New("models: email address is not valid")
	ErrEmailRequired     = errors.New("models: email address is required")
	ErrEmailTaken        = errors.New("models: email address is already taken")
	ErrPasswordTooShort  = errors.New("models: password must be at least 8 characters long")
	ErrPasswordRequired  = errors.New("models: password required")
	ErrRememberTooShort  = errors.New("models: remember token must be at least 32 bytes")
	ErrRememberRequired  = errors.New("models: remember token is required")
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
	Assets       []Asset
}

type userValFunc func(*User) error

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
	db *gorm.DB
}

// Interface fulfullment check
var _ UserDB = &userGorm{}

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
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
			hmac:   hash.NewHMAC(hmacSecretKey),
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
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

/////////////////////////////////// userValidator functions //////////////////////////////////

/*
func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
	}
}
*/

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setInitalRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.Nbytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

func (uv *userValidator) isGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
		return err
	}

	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

////////////////////////////////////////////// userValidator ////////////////////////////////////

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setInitalRemember,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.isGreaterThanZero)
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) requireEmail(user *User) error {
	fmt.Println("email: ", user.Email, "password: ", user.Password)
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

/////////////////////////////////// userGorm //////////////////////////////////////

func NewUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
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

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	db := ug.db.Where("remember_hash=?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
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
