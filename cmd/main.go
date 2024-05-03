package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		host     string
		port     int
		user     string
		password string
		dbname   string
	}
}

type application struct{}

func openDB(cfg config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.db.host, cfg.db.port, cfg.db.user, cfg.db.password, cfg.db.dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	cfg := config{
		port: 8080,
		env:  "development",
		db: struct {
			host     string
			port     int
			user     string
			password string
			dbname   string
		}{
			host:     "localhost",
			port:     5432,
			user:     "Aidyn",
			password: "050208551027",
			dbname:   "gun",
		},
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	log.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
