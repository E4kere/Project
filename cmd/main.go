package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/E4kere/Project/auth"
	"github.com/E4kere/Project/controller"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	auth.LoadEnvVariables()
	auth.ConnectToDB()
	auth.SyncDB()

}

type PaginatedResponse struct {
	TotalRecords int   `json:"totalRecords"`
	TotalPages   int   `json:"totalPages"`
	PageSize     int   `jsonjson:"pageSize"`
	CurrentPage  int   `json:"currentPage"`
	Data         []Gun `json:"data"`
}

type application struct {
	mu     sync.Mutex
	guns   map[int]Gun
	nextID int
	db     *sqlx.DB
}

// Gun struct
type Gun struct {
	ID     int
	Name   string
	Price  float64
	Damage int
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404: Not Found")
}
func ReadGun(db *sqlx.DB, id int) (*Gun, error) {
	var gun Gun
	err := db.Get(&gun, "SELECT id, name, price, damage FROM guns WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &gun, nil
}

func main() {
	// Database connection
	db, err := sqlx.Connect("postgres", "user=postgres dbname=gun sslmode=disable password=050208551027 host=localhost")
	if err != nil {
		log.Fatalln(err)
	}
	r := gin.Default()
	r.POST("/signup", controller.Signup)

	r.Run()

	defer db.Close()

	app := &application{db: db}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: app.routes(),
	}

	log.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
func (app *application) routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Protected routes
	s := r.PathPrefix("/").Subrouter()

	s.HandleFunc("/guns", app.listGuns).Methods("GET")
	s.HandleFunc("/guns/{id}", app.getGunByID).Methods("GET")
	s.HandleFunc("/guns", app.createGun).Methods("POST")
	s.HandleFunc("/guns/{id}", app.updateGun).Methods("PUT")
	s.HandleFunc("/guns/{id}", app.deleteGun).Methods("DELETE")

	return r
}

func (app *application) listGuns(w http.ResponseWriter, r *http.Request) {
	// Pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Filter parameters
	name := r.URL.Query().Get("name")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	minDamage := r.URL.Query().Get("minDamage")
	maxDamage := r.URL.Query().Get("maxDamage")

	// Construct SQL query dynamically
	query := "SELECT id, name, price, damage FROM guns WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Add filters dynamically
	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+name+"%")
		argIndex++
	}

	if minPrice != "" {
		minPriceVal, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			http.Error(w, "Invalid minPrice value", http.StatusBadRequest)
			return
		}
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, minPriceVal)
		argIndex++
	}

	if maxPrice != "" {
		maxPriceVal, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			http.Error(w, "Invalid maxPrice value", http.StatusBadRequest)
			return
		}
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, maxPriceVal)
		argIndex++
	}

	if minDamage != "" {
		minDamageVal, err := strconv.Atoi(minDamage)
		if err != nil {
			http.Error(w, "Invalid minDamage value", http.StatusBadRequest)
			return
		}
		query += fmt.Sprintf(" AND damage >= $%d", argIndex)
		args = append(args, minDamageVal)
		argIndex++
	}

	if maxDamage != "" {
		maxDamageVal, err := strconv.Atoi(maxDamage)
		if err != nil {
			http.Error(w, "Invalid maxDamage value", http.StatusBadRequest)
			return
		}
		query += fmt.Sprintf(" AND damage <= $%d", argIndex)
		args = append(args, maxDamageVal)
		argIndex++
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	// Execute the query
	var guns []Gun
	err = app.db.Select(&guns, query, args...)
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

	// Get total record count with filtering
	countQuery := "SELECT COUNT(*) FROM guns WHERE 1=1"
	if len(args) > 0 {
		countQuery += query[41:strings.Index(query, "LIMIT")]
	}
	err = app.db.Get(&response.TotalRecords, countQuery, args[:argIndex-1]...)
	if err != nil {
		http.Error(w, "Unable to count records", http.StatusInternalServerError)
		return
	}

	response.TotalPages = (response.TotalRecords + pageSize - 1) / pageSize

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) getGunByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var gun Gun
	err = app.db.Get(&gun, "SELECT id, name, price, damage FROM guns WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Gun not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) createGun(w http.ResponseWriter, r *http.Request) {
	var gun Gun
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

func (app *application) updateGun(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var gun Gun
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

/*func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Your login handler implementation here1
}
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the header is empty
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		token := tokenParts[1]

		// Perform token validation
		if token != "your_secret_token" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
*/
