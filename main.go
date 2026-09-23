package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BosBJJ/volunteer_app/database"
	"github.com/BosBJJ/volunteer_app/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	ctx := context.Background()
	conn, err := database.ConnectDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	err = database.CreateSchema(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", Welcome)
	http.HandleFunc("/volunteer", handlers.VolunteerHandler(conn))
	http.HandleFunc("/volunteer/{id}", handlers.GetVolunteerByID(conn))
	http.HandleFunc("/opportunity", handlers.OpportunityHandler(conn))
	http.HandleFunc("/opportunity/{id}", handlers.GetOpportunityByID(conn))
	http.HandleFunc("/opportunity/{id}/shifts", handlers.GetShiftsByOpportunityID(conn))
	http.HandleFunc("/shifts", handlers.ShiftHandler(conn))

	fmt.Println("Server started on port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}
