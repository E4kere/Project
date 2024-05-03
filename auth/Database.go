package auth

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDB() {
	var err error

	dsn := os.Getenv("DB")

	fmt.Println("DSN:", dsn)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Database connection error:", err)
		panic("Failed to connect to the database")
	}

}
