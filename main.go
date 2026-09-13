package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BosBJJ/volunteer_app/database"
	"github.com/BosBJJ/volunteer_app/handlers"
)

func main() {
	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	err = database.CreateSchema(conn)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", Welcome)
	http.HandleFunc("/volunteer", handlers.VolunteerHandler(conn))
	http.HandleFunc("/volunteer/{id}", handlers.GetVolunteerByID(conn))
	http.HandleFunc("/opportunity", handlers.OpportunityHandler(conn))
	http.HandleFunc("/opportunity/{id}", handlers.GetOpportunityByID(conn))

	fmt.Println("Server started on port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}
