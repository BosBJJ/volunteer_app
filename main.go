package main

import (
	"fmt"
	"net/http"

	"github.com/BosBJJ/volunteer_app/handlers"
)

func main() {
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/volunteer", handlers.VolunteerHandler)
	http.HandleFunc("/volunteer/{id}", handlers.GetVolunteerByID)
	http.HandleFunc("/opportunity", handlers.OpportunityHandler)

	fmt.Println("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}
