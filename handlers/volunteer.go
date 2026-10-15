package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func VolunteerByIDHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			GetVolunteerByID(conn, w, req)
		case http.MethodDelete:
			DeleteVolunteerByID(conn, w, req)
		case http.MethodPatch:
			UpdateVolunteerByID(conn, w, req)
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

func GetVolunteerByID(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
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

func ShowAttendance(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		volunteerID, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid volunteer id", http.StatusBadRequest)
			return
		}
		volExists, err := database.VolunteerExists(req.Context(), conn, volunteerID)
		if err != nil {
			http.Error(w, "unable to verify volunteer", http.StatusInternalServerError)
			return
		}
		if !volExists {
			http.Error(w, "volunteer not found", http.StatusNotFound)
			return
		}
		records, err := database.GetAttendanceByVolunteerID(req.Context(), conn, volunteerID)
		if err != nil {
			http.Error(w, "error fetching attendance", http.StatusInternalServerError)
			return
		}
		hours, err := database.GetTotalHoursByVolunteerID(req.Context(), conn, volunteerID)
		if err != nil {
			http.Error(w, "error fetching hours", http.StatusInternalServerError)
			return
		}
		resp := models.AttendanceResponse{
			Records:    records,
			TotalHours: hours,
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func DeleteVolunteerByID(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	volunteerID, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid volunteer id", http.StatusBadRequest)
		return
	}
	err = database.DeleteVolunteer(req.Context(), conn, volunteerID)
	if err != nil {
		if errors.Is(err, database.ErrVolunteerNotFound) {
			http.Error(w, "volunteer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting volunteer", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateVolunteerByID(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	volunteerID, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid volunteer id", http.StatusBadRequest)
		return
	}
	var volunteer models.VolunteerUpdate
	err = json.NewDecoder(req.Body).Decode(&volunteer)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if volunteer.Name == nil && volunteer.Email == nil {
		http.Error(w, "at least one field must be provided", http.StatusBadRequest)
		return
	}
	if volunteer.Name != nil {
		*volunteer.Name = strings.TrimSpace(*volunteer.Name)
	}
	if volunteer.Email != nil {
		*volunteer.Email = strings.TrimSpace(*volunteer.Email)
	}
	if (volunteer.Name != nil && *volunteer.Name == "") || (volunteer.Email != nil && *volunteer.Email == "") {
		http.Error(w, "name and email can't be empty", http.StatusBadRequest)
		return
	}
	err = database.EditVolunteer(req.Context(), conn, volunteerID, volunteer)
	if err != nil {
		if errors.Is(err, database.ErrVolunteerNotFound) {
			http.Error(w, "volunteer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating volunteer", http.StatusInternalServerError)
		return
	}
	updatedVolunteer, err := database.ListVolunteerByID(req.Context(), conn, volunteerID)
	if err != nil {
		http.Error(w, "error fetching updated  volunteer", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(updatedVolunteer)
}
