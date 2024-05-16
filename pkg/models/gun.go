package models

import (
	"database/sql"
	"log"
	"time"
)

type Gun struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Damage    int       `json:"damage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GunModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
