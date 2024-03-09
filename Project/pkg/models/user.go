package project

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

type User struct {
	Id       int
	Username string
	Password string
}

func (m *UserModel) RegisterUser(login string, password string, writer http.ResponseWriter) error {
	var username string
	err := m.DB.QueryRow("SELECT username FROM users WHERE username = $1", login).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			hashPass, err := HashPassword(password)
			if err != nil {
				http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
				return err
			}
			_, err = m.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", login, hashPass)
			if err != nil {
				http.Error(writer, "Failed to insert user into database", http.StatusInternalServerError)
				return err
			}
			_, err = fmt.Fprintf(writer, "User registered successfully")
			if err != nil {
				return err
			}
			return nil
		}
		http.Error(writer, "Error checking existing user", http.StatusInternalServerError)
		return err
	}

	http.Error(writer, "User already exists", http.StatusBadRequest)
	return nil
}
func (m *UserModel) LoginUser(login string, password string, writer http.ResponseWriter) error {
	rows := m.DB.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)
	if user.Username == "" {
		_, err := fmt.Fprintf(writer, "User not found")
		if err != nil {
			return err
		}
		return nil
	}
	if CheckPasswordHash(password, user.Password) {
		payload := jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		t, err := token.SignedString([]byte(os.Getenv("secretKey")))
		fmt.Println(os.Getenv("secretKey"))

		if err != nil {
			http.Error(writer, "jwt token signing", http.StatusBadRequest)
		}

		err = json.NewEncoder(writer).Encode(t)
		if err != nil {
			return err
		}

	} else {
		_, err := fmt.Fprintf(writer, "Wrong password")
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *UserModel) DeleteUser(login string, writer http.ResponseWriter) {
	rows := m.DB.QueryRow("SELECT * FROM users WHERE username = $1", login)
	user := User{}
	_ = rows.Scan(&user.Id, &user.Username, &user.Password)

	if rows != nil {
		_, err := m.DB.Exec("DELETE FROM users WHERE username = $1", login)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(writer, "User deleted")
		if err != nil {
			return
		}
	} else {
		_, err := fmt.Fprintf(writer, "User not found write credentials")
		if err != nil {
			return
		}
	}
}

func (m *UserModel) GetAllUsers(writer http.ResponseWriter) {
	rows, err := m.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(writer, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	var users []User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		user.Password = ""
		fmt.Println(user)
		if err != nil {
			http.Error(writer, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	_, err = fmt.Fprintf(writer, "%v", users)
	if err != nil {
		return
	}

}
func setupRoutes(router *mux.Router, userModel *UserModel) {
	router.HandleFunc("/players", GetPlayersHandler(userModel)).Methods("GET")

}

func GetPlayersHandler(userModel *UserModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
