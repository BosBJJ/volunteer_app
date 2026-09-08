package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/BosBJJ/volunteer_app/models"
)

var volunteers []models.Volunteer

func VolunteerHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		CreateVolunteer(w, req)
	case http.MethodGet:
		GetVolunteers(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// temporary until SQL implemented
var nextID int = 0

func CreateVolunteer(w http.ResponseWriter, req *http.Request) {
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
	volunteer.RegisteredAt = time.Now()
	volunteer.Id = nextID
	volunteers = append(volunteers, volunteer)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volunteer)
	nextID++
}

func GetVolunteers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volunteers)
}

func GetVolunteerByID(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("id")
	reqID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	for _, volunteer := range volunteers {
		if volunteer.Id == reqID {
			w.Header().Set("content-type", "application/json")
			json.NewEncoder(w).Encode(volunteer)
			return
		}
	}
	http.Error(w, "volunteer not found", http.StatusNotFound)
}
