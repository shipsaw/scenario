package controllers

import (
	"fmt"
	"net/http"

	"github.com/shipsaw/scenario/models"
	"github.com/shipsaw/scenario/views"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Used to create a new Users controller.  Should
// only be used during initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
		us:      us,
	}
}

// Used to render the form where user can create new user account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Used to process the signup form when a user
// submits it.  Creates a new user account
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
