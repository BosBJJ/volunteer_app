package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var volunteers []Volunteer

func main() {
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/volunteer", CreateVolunteer)

	fmt.Println("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}

func CreateVolunteer(w http.ResponseWriter, req *http.Request) {
	var volunteer Volunteer
	err := json.NewDecoder(req.Body).Decode(&volunteer)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	volunteer.RegisteredAt = time.Now()
	volunteers = append(volunteers, volunteer)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(volunteer)
}

type Volunteer struct {
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registeredAt"`
}
