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
	galleriesC := controllers.NewGalleries(services.Gallery, router)

	router.Handle("/", staticC.HomeView).Methods("GET")
	router.Handle("/contact", staticC.ContactView).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	router.Handle("/login", usersC.LoginView).Methods("GET")
	router.HandleFunc("/login", usersC.Login).Methods("POST")
	router.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// Gallery routes
	requireUserMw := middleware.RequireUser{UserService: services.User}
	router.Handle("/galleries/new", requireUserMw.Apply(galleriesC.NewView)).Methods("GET")
	router.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET")
	router.HandleFunc("/galleries/:id", galleriesC.Show).Methods("GET")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	http.ListenAndServe(":3000", router)
}
