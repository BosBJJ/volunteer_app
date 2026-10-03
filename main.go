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
	http.HandleFunc("/volunteer/{id}/attendance", handlers.ShowAttendance(conn))
	http.HandleFunc("/opportunity", handlers.OpportunityHandler(conn))
	http.HandleFunc("/opportunity/{id}", handlers.GetOpportunityByID(conn))
	http.HandleFunc("/opportunity/{id}/shifts", handlers.GetShiftsByOpportunityID(conn))
	http.HandleFunc("/shifts", handlers.ShiftHandler(conn))
	http.HandleFunc("/shifts/{id}/signup", handlers.SignupToShift(conn)) //temporary use - http://localhost:8080/shifts/1/signup?volunteerId=1
	http.HandleFunc("/shifts/{id}/signups", handlers.CheckSignupsByShift(conn))
	http.HandleFunc("/shifts/{id}/checkin", handlers.ShiftCheckIn(conn))   //temporary use - http://localhost:8080/shifts/1/checkin?volunteerId=1
	http.HandleFunc("/shifts/{id}/checkout", handlers.ShiftCheckOut(conn)) //temporary use - http://localhost:8080/shifts/1/checkout?volunteerId=1

	fmt.Println("Server started on port 8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}
