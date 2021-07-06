package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shipsaw/scenario/controllers"
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
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()
	// us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	router := mux.NewRouter()
	router.Handle("/", staticC.HomeView).Methods("GET")
	router.Handle("/contact", staticC.ContactView).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", router)
}
