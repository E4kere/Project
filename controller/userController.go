package controller

import (
	"net/http"

	"github.com/E4kere/Project/auth"
	"github.com/E4kere/Project/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read",
		})

		return
	}

	//Hash
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error hashing",
		})

		return
	}

	//Create
	user := models.User{Email: body.Email, Password: string(hash)}

	result := auth.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already exists",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{})

}
