package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/E4kere/Project/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type application struct {
	db       *sqlx.DB
	errorLog *log.Logger
	models   models.Models
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

	// Set up routes
	router := mux.NewRouter()
	// router.HandleFunc("/register", app.Register).Methods("POST")
	// router.HandleFunc("/login", app.Login).Methods("POST")

	// router.Use(app.Authenticate)

	// CRUD routes for guns
	gunsRouter := router.PathPrefix("/guns").Subrouter()
	// gunsRouter.Use(app.Authenticate)
	gunsRouter.HandleFunc("", app.listGuns).Methods("GET")
	gunsRouter.HandleFunc("", app.createGun).Methods("POST")
	gunsRouter.HandleFunc("/{id}", app.getGunByID).Methods("GET")
	gunsRouter.HandleFunc("/{id}", app.updateGun).Methods("PUT")
	gunsRouter.HandleFunc("/{id}", app.deleteGun).Methods("DELETE")

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

// func (app *application) Register(w http.ResponseWriter, r *http.Request) {
// 	var user struct {
// 		Name     string `json:"name"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	query := `
// 			INSERT INTO users (name, email, password_hash)
// 			VALUES ($1, $2, $3)
// 			RETURNING id
// 	`
// 	var userID int64
// 	err = app.db.QueryRow(query, user.Name, user.Email, hashedPassword).Scan(&userID)
// 	if err != nil {
// 		http.Error(w, "Unable to create user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Generate JWT token
// 	tokenString, err := app.generateToken(userID)
// 	if err != nil {
// 		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set JWT token as cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "token",
// 		Value:    tokenString,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: true,
// 		SameSite: http.SameSiteStrictMode,
// 	})

// 	// Return success response
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "User registered successfully",
// 		"token":   tokenString,
// 	})
// }

// func (app *application) Login(w http.ResponseWriter, r *http.Request) {
// 	var user struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	var userID int64
// 	var hashedPassword []byte
// 	err := app.db.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", user.Email).Scan(&userID, &hashedPassword)
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password)); err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	// Generate JWT token
// 	tokenString, err := app.generateToken(userID)
// 	if err != nil {
// 		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set JWT token as cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "token",
// 		Value:    tokenString,
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: true,
// 		SameSite: http.SameSiteStrictMode,
// 	})

// 	// Return success response
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Login successful",
// 		"token":   tokenString,
// 	})
// }

// func (app *application) Authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("token")
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := cookie.Value
// 		claims := &Claims{}

// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return app.jwtSecret, nil
// 		})
// 		if err != nil || !token.Valid {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "user_id", strconv.FormatInt(claims.UserID, 10))

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func (app *application) generateToken(userID int64) (string, error) {
// 	claims := &Claims{
// 		UserID: userID,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
// 			Issuer:    "your-issuer",
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(app.jwtSecret)
// }

// // Automated tests in Postman
// func runTestsInPostman() {
// 	// Test GET request to retrieve data
// 	testGetRequest("Retrieve Data", "http://localhost:8080/guns")

// 	// Test POST request to create data
// 	testPostRequest("Create Data", "http://localhost:8080/guns", map[string]interface{}{
// 		"name":   "Gun Name",
// 		"price":  100,
// 		"damage": 50,
// 	})

// 	// Test DELETE request to delete data

// 	testDeleteRequest("Delete Data", "http://localhost:8080/guns/100")

// }

// func testGetRequest(testName, url string) {
// 	log.Printf("Running test: %s", testName)
// 	response, err := http.Get(url)
// 	if err != nil {
// 		log.Fatalf("Error executing GET request: %v", err)
// 	}
// 	defer response.Body.Close()

// 	log.Printf("Response status: %s", response.Status)

// }

// func testPostRequest(testName, url string, data interface{}) {
// 	log.Printf("Running test: %s", testName)
// 	payload, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Error marshalling data: %v", err)
// 	}

// 	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
// 	if err != nil {
// 		log.Fatalf("Error executing POST request: %v", err)
// 	}
// 	defer response.Body.Close()

// 	log.Printf("Response status: %s", response.Status)
// 	// Additional assertions can be added here to verify the response
// }

// func testDeleteRequest(testName, url string) {
// 	log.Printf("Running test: %s", testName)
// 	client := &http.Client{}
// 	req, err := http.NewRequest("DELETE", url, nil)
// 	if err != nil {
// 		log.Fatalf("Error creating DELETE request: %v", err)
// 	}

// 	response, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Error executing DELETE request: %v", err)
// 	}
// 	defer response.Body.Close()

// 	log.Printf("Response status: %s", response.Status)

// }
