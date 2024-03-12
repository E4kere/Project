package models

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Guns GunModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\tCS2Guns\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\tCS2Guns\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Guns: GunModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
