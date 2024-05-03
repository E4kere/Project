package auth

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {

	var err error
	dsn := "host=localhost user=postgres password=050208551027 dbname=auth port=5432 sslmode=disable "

	fmt.Println("DSN:", dsn)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Database connection error:", err)
		panic("Failed to connect to the database")
	}

}
