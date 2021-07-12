package models

import (
	"strings"
)

const (
	// User-related errors
	ErrNotFound          modelError   = "models: resource not found"
	ErrPasswordIncorrect modelError   = "models: incorrect password provided"
	ErrPasswordTooShort  modelError   = "models: password must be at least 8 characters long"
	ErrEmailInvalid      modelError   = "models: email address is not valid"
	ErrEmailRequired     modelError   = "models: email address is required"
	ErrEmailTaken        modelError   = "models: email address is already taken"
	ErrPasswordRequired  modelError   = "models: password required"
	ErrRememberTooShort  privateError = "models: remember token must be at least 32 bytes"
	ErrRememberRequired  privateError = "models: remember token is required"
	ErrIDInvalid         privateError = "models: id provided is invalid"

	// Gallery-related errors
	ErrUserIDRequired privateError = "models: user ID is required"
	ErrTitleRequired  modelError   = "models: title is required"
)

type modelError string

type privateError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

func (e privateError) Error() string {
	return string(e)
}
