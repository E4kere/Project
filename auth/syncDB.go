package auth

import "github.com/E4kere/Project/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
