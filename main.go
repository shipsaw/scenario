package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shipsaw/scenario/controllers"
	"github.com/shipsaw/scenario/middleware"
	"github.com/shipsaw/scenario/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	// services.DestructiveReset()
	services.AutoMigrate()

	router := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, router)

	router.Handle("/", staticC.HomeView).Methods("GET")
	router.Handle("/contact", staticC.ContactView).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	router.Handle("/login", usersC.LoginView).Methods("GET")
	router.HandleFunc("/login", usersC.Login).Methods("POST")

	// Gallery routes
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}
	router.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	router.Handle("/galleries/new", requireUserMw.Apply(galleriesC.NewView)).Methods("GET")
	router.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	router.HandleFunc("/galleries/:id", galleriesC.Show).Methods("GET")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

	http.ListenAndServe(":3000", userMw.Apply(router))
}
