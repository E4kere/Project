package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	r.HandleFunc("/guns", app.listGuns).Methods("GET")
	r.HandleFunc("/guns/{id}", app.getGunByID).Methods("GET")
	r.HandleFunc("/guns", app.createGun).Methods("POST")
	r.HandleFunc("/guns/{id}", app.updateGun).Methods("PUT")
	r.HandleFunc("/guns/{id}", app.deleteGun).Methods("DELETE")

	return r
}
