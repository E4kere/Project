package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/Zhan1bek/BookStore/pkg/models"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models models.Models
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createGunHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string  `json:"name"`
		Manufacturer string  `json:"manufacturer"`
		Price        float64 `json:"price"`
		Damage       int     `json:"damage"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	gun := &models.Gun{
		Name:         input.Name,
		Manufacturer: input.Manufacturer,
		Price:        input.Price,
		Damage:       input.Damage,
	}

	err = app.models.Guns.Insert(gun)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, gun)
}

func (app *application) getGunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid gun ID")
		return
	}

	gun, err := app.models.Guns.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Gun not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, gun)
}

func (app *application) deleteGunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid gun ID")
		return
	}

	err = app.models.Guns.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) updateGunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid gun ID")
		return
	}

	var input struct {
		Name         *string  `json:"name"`
		Manufacturer *string  `json:"manufacturer"`
		Price        *float64 `json:"price"`
		Damage       *int     `json:"damage"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	gun := models.Gun{ID: int64(id)}
	if input.Name != nil {
		gun.Name = *input.Name
	}
	if input.Manufacturer != nil {
		gun.Manufacturer = *input.Manufacturer
	}
	if input.Price != nil {
		gun.Price = *input.Price
	}
	if input.Damage != nil {
		gun.Damage = *input.Damage
	}

	err = app.models.Guns.Update(&gun)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, gun)
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgresql://Aidyn:050208551027@localhost/gunstore?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: models.NewModels(db),
	}

	app.run()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) run() {
	r := mux.NewRouter()

	// Setting up routes for CS2 guns
	gunRouter := r.PathPrefix("/api/v1/guns").Subrouter()
	gunRouter.HandleFunc("", app.createGunHandler).Methods("POST")               // Create a gun
	gunRouter.HandleFunc("/{id:[0-9]+}", app.getGunHandler).Methods("GET")       // Get a gun by ID
	gunRouter.HandleFunc("/{id:[0-9]+}", app.updateGunHandler).Methods("PUT")    // Update a gun
	gunRouter.HandleFunc("/{id:[0-9]+}", app.deleteGunHandler).Methods("DELETE") // Delete a gun

	// Starting the server
	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	if err != nil {
		log.Fatal(err)
	}
}
