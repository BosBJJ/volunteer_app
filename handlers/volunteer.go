package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BosBJJ/volunteer_app/database"
	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

func VolunteerHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			CreateVolunteer(conn, w, req)
		case http.MethodGet:
			GetVolunteers(conn, w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func CreateVolunteer(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	var volunteer models.Volunteer
	err := json.NewDecoder(req.Body).Decode(&volunteer)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if volunteer.Name == "" || volunteer.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}
	err = database.SaveVolunteer(req.Context(), conn, &volunteer)
	if err != nil {
		http.Error(w, "error saving volunteer to database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volunteer)
}

func GetVolunteers(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	volunteers, err := database.ListVolunteers(req.Context(), conn)
	if err != nil {
		http.Error(w, "error fetching volunteers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volunteers)
}

func GetVolunteerByID(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.PathValue("id")
		reqID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		volunteer, err := database.ListVolunteerByID(req.Context(), conn, reqID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "invalid volunteer id", http.StatusNotFound)
			} else {
				http.Error(w, "error fetching volunteer", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(volunteer)
	}
}
