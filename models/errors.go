package models

import (
	"errors"
	"strings"
)

const (
	ErrNotFound          modelError = "models: resource not found"
	ErrIDInvalid         modelError = "models: id provided is invalid"
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	ErrPasswordTooShort  modelError = "models: password must be at least 8 characters long"
	ErrEmailInvalid      modelError = "models: email address is not valid"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailTaken        modelError = "models: email address is already taken"
	ErrPasswordRequired  modelError = "models: password required"
)

var (
	ErrRememberTooShort = errors.New("models: remember token must be at least 32 bytes")
	ErrRememberRequired = errors.New("models: remember token is required")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
