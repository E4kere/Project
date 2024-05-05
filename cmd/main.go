package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/E4kere/Project/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type application struct {
	db     *sqlx.DB
	models models.Models
	jwtKey []byte
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to the database using sqlx
	db, err := sqlx.Connect("postgres", "postgres://postgres:050208551027@localhost:5432/gun?sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Initialize the application struct
	app := &application{
		db:     db,
		models: models.NewModels(db),
		jwtKey: []byte("be5906c88cb6519f83fadf949bb2cde051b288613e73845182f473c911464af3")}

	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/register", app.Register).Methods("POST")
	router.HandleFunc("/login", app.Login).Methods("POST")

	// CRUD routes for guns
	router.HandleFunc("/guns", app.listGuns).Methods("GET")
	router.HandleFunc("/guns", app.createGun).Methods("POST")
	router.HandleFunc("/guns/{id}", app.getGunByID).Methods("GET")
	router.HandleFunc("/guns/{id}", app.updateGun).Methods("PUT")
	router.HandleFunc("/guns/{id}", app.deleteGun).Methods("DELETE")

	// Start the server
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (app *application) listGuns(w http.ResponseWriter, r *http.Request) {
	var guns []models.Gun
	err := app.db.Select(&guns, "SELECT id, name, price, damage FROM guns ORDER BY id")
	if err != nil {
		http.Error(w, "Unable to fetch guns", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guns)
}

func (app *application) createGun(w http.ResponseWriter, r *http.Request) {
	var gun models.Gun
	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO guns (name, price, damage) VALUES ($1, $2, $3) RETURNING id"
	err := app.db.QueryRow(query, gun.Name, gun.Price, gun.Damage).Scan(&gun.ID)
	if err != nil {
		http.Error(w, "Unable to create gun", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) getGunByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var gun models.Gun
	err = app.db.Get(&gun, "SELECT id, name, price, damage FROM guns WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Gun not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) updateGun(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var gun models.Gun
	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	gun.ID = id
	query := "UPDATE guns SET name = $1, price = $2, damage = $3 WHERE id = $4"
	_, err = app.db.Exec(query, gun.Name, gun.Price, gun.Damage, gun.ID)
	if err != nil {
		http.Error(w, "Unable to update gun", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) deleteGun(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM guns WHERE id = $1"
	_, err = app.db.Exec(query, id)
	if err != nil {
		http.Error(w, "Unable to delete gun", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var credentials models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}

	// Call Insert on Users model
	err = app.models.Users.Insert("User Name", credentials.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := app.models.Users.GetByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Use Matches to check the password
	if !user.Password.Matches(credentials.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := app.models.Tokens.New(user.ID, app.jwtKey)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
}
