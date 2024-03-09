package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Player struct {
	Name        string
	Team        string
	Country     string
	DateOfBirth float64
}

func main() {
	connStr := "postgres://Aidyn:050208551027@localhost:5432/gopgtest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	CreateTablePlayer(db)

	// Create a new Router
	router := mux.NewRouter()

	// Initialize User model with DB connection
	userModel := models.NewUserModel(db)
	// Setup routes for the API
	setupRoutes(router, userModel)

	// Start the server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func CreateTablePlayer(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS Players (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			team VARCHAR(255),
			country VARCHAR(255),
			date_of_birth DATE
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
