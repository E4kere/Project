package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "OK")
	if err != nil {
		return
	}
}
func RegisterHandler(userModel *UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		login := request.FormValue("login")
		password := request.FormValue("password")

		err := userModel.RegisterUser(login, password, writer)
		if err != nil {

			return
		}
	}
}

func Login(userModel *UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		login := request.FormValue("login")
		password := request.FormValue("password")

		err := userModel.LoginUser(login, password, writer)
		if err != nil {
			return
		}
	}
}
func Update(userModel *UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		updateType := mux.Vars(request)["type"]
		login := payload["sub"].(string)
		userModel.UpdateUser(login, updateType, writer, request)
	}
}

func DeleteUserHandler(userModel *UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		if !check {
			return
		}
		login := payload["sub"].(string)
		userModel.DeleteUser(login, writer)
	}
}

func GetAllUsersHandler(userModel *UserModel) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload, check := JwtPayloadFromRequest(writer, request)
		fmt.Println(payload["sub"])
		if !check {
			return
		}
		userModel.GetAllUsers(writer)
	}
}
func SendMessage(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "SendMessage")
	if err != nil {
		return
	}
}
func UpdateMessage(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "UpdateMessage")
	if err != nil {
		return
	}
}
func DeleteMessage(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "DeleteMessage")
	if err != nil {
		return
	}
}
func Notifications(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprintf(writer, "Notifications")
	if err != nil {
		return
	}
}
