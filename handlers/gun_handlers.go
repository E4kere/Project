package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codev0/inft3212-6/pkg/abr-plus/model"
	"github.com/gorilla/mux"
)

func (app *application) listGuns(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Logic to fetch paginated list of guns (dummy data for now)
	guns := []model.Gun{
		{ID: 1, Name: "Gun A", Type: "Handgun", Manufacturer: "Company A"},
		{ID: 2, Name: "Gun B", Type: "Rifle", Manufacturer: "Company B"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guns)
}

func (app *application) getGunByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid gun ID", http.StatusBadRequest)
		return
	}

	// Logic to fetch a single gun by ID (dummy data for now)
	gun := model.Gun{ID: id, Name: "Gun A", Type: "Handgun", Manufacturer: "Company A"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) createGun(w http.ResponseWriter, r *http.Request) {
	var gun model.Gun
	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Logic to save a new gun to the database (dummy data for now)
	gun.ID = 3 // Dummy ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) updateGun(w httpResponseWriter, r *HttpRequest) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid gun ID", http.StatusBadRequest)
		return
	}

	var gun model.Gun
	if err := json.NewDecoder(r.Body).Decode(&gun); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	gun.ID = id // Dummy update
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gun)
}

func (app *application) deleteGun(w httpResponseWriter, r *HttpRequest) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid gun ID", http.StatusBadRequest)
		return
	}

	// Logic to delete a gun by ID (dummy data for now)
	w.WriteHeader(http.StatusNoContent)
}
