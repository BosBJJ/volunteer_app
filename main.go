package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var volunteers []Volunteer

func main() {
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/volunteer", VolunteerHandler)
	http.HandleFunc("/volunteer/{id}", GetVolunteerByID)

	fmt.Println("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}

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
	var volunteer Volunteer
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
func GetVolunteers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volunteers)
}

type Volunteer struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registeredAt"`
}
