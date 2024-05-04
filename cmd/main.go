package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/E4kere/Project/auth"
	"github.com/E4kere/Project/controller"
	"github.com/E4kere/Project/middleware"
	"github.com/E4kere/Project/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

func main() {
	auth.ConnectToDB()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.HandleFunc("/signup", controller.Signup).Methods("POST")

	// Protected routes
	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)
	protectedRouter.HandleFunc("/guns", listGuns).Methods("GET")
	protectedRouter.HandleFunc("/guns", createGun).Methods("POST")
	protectedRouter.HandleFunc("/guns/{id}", deleteGun).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handler to list guns
func listGuns(w http.ResponseWriter, r *http.Request) {
	var guns []models.Gun
	err := db.Select(&guns, "SELECT id, name, price, damage FROM guns")
	if err != nil {
		http.Error(w, "Unable to fetch guns", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guns)
}

// Handler to create a gun
func createGun(w http.ResponseWriter, r *http.Request) {
	var gun models.Gun
	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO guns (name, price, damage) VALUES ($1, $2, $3) RETURNING id"
	err := db.QueryRow(query, gun.Name, gun.Price, gun.Damage).Scan(&gun.ID)
	if err != nil {
		http.Error(w, "Unable to create gun", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

// Handler to delete a gun
func deleteGun(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM guns WHERE id = $1"
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, "Unable to delete gun", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"sync"

// 	"github.com/E4kere/Project/models"

// 	"github.com/gorilla/mux"
// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

// type PaginatedResponse struct {
// 	TotalRecords int          `json:"totalRecords"`
// 	TotalPages   int          `json:"totalPages"`
// 	PageSize     int          `json:"pageSize"`
// 	CurrentPage  int          `json:"currentPage"`
// 	Data         []models.Gun `json:"data"`
// }

// type application struct {
// 	mu     sync.Mutex
// 	guns   map[int]models.Gun
// 	nextID int
// 	db     *sqlx.DB
// }

// func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusNotFound)
// 	fmt.Fprintf(w, "404: Not Found")
// }

// func main() {
// 	db, err := sqlx.Connect("postgres", "user=postgres dbname=gun sslmode=disable password=050208551027 host=localhost")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer db.Close()

// 	app := &application{db: db}
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%d", 8080),
// 		Handler: app.routes(),
// 	}

// 	log.Printf("Starting server on %s", srv.Addr)
// 	err = srv.ListenAndServe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func (app *application) routes() http.Handler {
// 	router.Run(":8080")

// func (app *application) listGuns(w http.ResponseWriter, r *http.Request) {
// 	// Get sorting parameters from the query string
// 	sortField := r.URL.Query().Get("sort")
// 	sortOrder := r.URL.Query().Get("order")

// 	// Default sorting by ID if not specified
// 	if sortField == "" {
// 		sortField = "id"
// 	}

// 	// Default order is ascending if not specified
// 	if sortOrder != "asc" && sortOrder != "desc" {
// 		sortOrder = "asc"
// 		}}}

// 	// Validate sorting field
// 	validSortFields := map[string]bool{
// 		"id": true, "name": true, "price": true, "damage": true,
// 	}
// 	if !validSortFields[sortField] {
// 		http.Error(w, "Invalid sort field", http.StatusBadRequest)
// 		return
// 	}

// 	// Construct the SQL query with sorting
// 	query := fmt.Sprintf("SELECT id, name, price, damage FROM guns ORDER BY %s %s", sortField, sortOrder)

// 	// Execute the query to fetch all guns
// 	var guns []models.Gun
// 	err := app.db.Select(&guns, query)
// 	if err != nil {
// 		http.Error(w, "Unable to retrieve guns", http.StatusInternalServerError)
// 		return
// 	}

// 	// Prepare the response
// 	response := PaginatedResponse{
// 		Data: guns,
// 	}

// 	// Calculate the total number of records (without pagination)
// 	response.TotalRecords = len(guns)
// 	response.TotalPages = 1
// 	response.PageSize = response.TotalRecords
// 	response.CurrentPage = 1

// 	// Encode and write the response as JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func (app *application) getGunByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := mux.Vars(r)["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var gun models.Gun
// 	err = app.db.Get(&gun, "SELECT id, name, price, damage FROM guns WHERE id = $1", id)
// 	if err != nil {
// 		http.Error(w, "Gun not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(gun)
// }

// func (app *application) createGun(w http.ResponseWriter, r *http.Request) {
// 	var gun models.Gun
// 	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	query := "INSERT INTO guns (name, price, damage) VALUES ($1, $2, $3) RETURNING id"
// 	err := app.db.QueryRow(query, gun.Name, gun.Price, gun.Damage).Scan(&gun.ID)
// 	if err != nil {
// 		http.Error(w, "Unable to create gun", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(gun)
// }

// func (app *application) updateGun(w http.ResponseWriter, r *http.Request) {
// 	idStr := mux.Vars(r)["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var gun models.Gun
// 	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	gun.ID = id
// 	query := "UPDATE guns SET name = $1, price = $2, damage = $3 WHERE id = $4"
// 	_, err = app.db.Exec(query, gun.Name, gun.Price, gun.Damage, gun.ID)
// 	if err != nil {
// 		http.Error(w, "Unable to update gun", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(gun)
// }

// func (app *application) deleteGun(w http.ResponseWriter, r *http.Request) {
// 	idStr := mux.Vars(r)["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	query := "DELETE FROM guns WHERE id = $1"
// 	_, err = app.db.Exec(query, id)
// 	if err != nil {
// 		http.Error(w, "Unable to delete gun", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
