package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/E4kere/Project/pkg/jsonlog"
	"github.com/E4kere/Project/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type application struct {
	db *sqlx.DB

	logger *jsonlog.Logger
	models models.Models
}

type PaginatedResponse struct {
	TotalRecords int         `json:"totalRecords"`
	TotalPages   int         `json:"totalPages"`
	PageSize     int         `json:"pageSize"`
	CurrentPage  int         `json:"currentPage"`
	Data         interface{} `json:"data"`
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
		db: db,
	}

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// gunsRouter.Use(app.Authenticate)
	router.HandleFunc("/guns", app.listGuns).Methods("GET")
	router.HandleFunc("/guns", app.createGun).Methods("POST")
	router.HandleFunc("/guns/{id}", app.getGunByID).Methods("GET")
	router.HandleFunc("/guns/{id}", app.updateGun).Methods("PUT")
	router.HandleFunc("/guns/{id}", app.deleteGun).Methods("DELETE")
	router.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	// Start the server
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
func (app *application) listGuns(w http.ResponseWriter, r *http.Request) {
	// Extract pagination parameters from the query string
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get sorting parameters from the query string
	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// Default sorting by ID if not specified
	if sortField == "" {
		sortField = "id"
	}

	// Default order is ascending if not specified
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Validate sorting field
	validSortFields := map[string]bool{
		"id": true, "name": true, "price": true, "damage": true,
	}
	if !validSortFields[sortField] {
		http.Error(w, "Invalid sort field", http.StatusBadRequest)
		return
	}

	// Construct the SQL query with sorting and pagination
	query := fmt.Sprintf(
		"SELECT id, name, price, damage FROM guns ORDER BY %s %s LIMIT $1 OFFSET $2",
		sortField, sortOrder)

	// Execute the query to fetch guns with pagination
	var guns []models.Gun
	err = app.db.Select(&guns, query, pageSize, offset)
	if err != nil {
		http.Error(w, "Unable to retrieve guns", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := PaginatedResponse{
		PageSize:    pageSize,
		CurrentPage: page,
		Data:        guns,
	}

	// Calculate the total number of records (without filtering)
	countQuery := "SELECT COUNT(*) FROM guns"
	err = app.db.Get(&response.TotalRecords, countQuery)
	if err != nil {
		http.Error(w, "Unable to count records", http.StatusInternalServerError)
		return
	}

	// Calculate the total number of pages
	response.TotalPages = (response.TotalRecords + pageSize - 1) / pageSize

	// Encode and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
