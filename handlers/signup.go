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

func SignupToShift(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		shiftID, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shift id", http.StatusBadRequest)
			return
		}
		//REMOVE WHEN AUTHORIZATION IMPLEMENTED
		volunteerID, err := strconv.Atoi(req.URL.Query().Get("volunteerId"))
		if err != nil {
			http.Error(w, "invalid volunteer id", http.StatusBadRequest)
			return
		}
		var signup models.Signup
		shift, err := database.ListShiftByShiftId(req.Context(), conn, shiftID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "shift not found", http.StatusNotFound)
				return
			}
			http.Error(w, "unable to fetch shift", http.StatusInternalServerError)
			return
		}
		volExists, err := database.VolunteerExists(req.Context(), conn, volunteerID)
		if err != nil {
			http.Error(w, "unable to verify volunteer", http.StatusInternalServerError)
			return
		}
		if !volExists {
			http.Error(w, "invalid volunteer id", http.StatusBadRequest)
			return
		}
		exists, err := database.SignupExists(req.Context(), conn, volunteerID, shiftID)
		if err != nil {
			http.Error(w, "unable to fetch signups", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "user already registered for this shift", http.StatusConflict)
			return
		}
		signup.ShiftID = shiftID
		signup.VolunteerID = volunteerID
		registerCount, err := database.CountSignups(req.Context(), conn, shiftID)
		if err != nil {
			http.Error(w, "unable to fetch signups", http.StatusInternalServerError)
			return
		}
		if registerCount >= shift.Capacity {
			http.Error(w, "shift is full", http.StatusConflict)
			return
		}
		err = database.SaveSignup(req.Context(), conn, &signup)
		if err != nil {
			http.Error(w, "unable to sign up", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(signup)
	}
}

func CheckSignupsByShift(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		shiftID, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid shift id", http.StatusBadRequest)
			return
		}
		volunteers, err := database.ListSignupsByShift(req.Context(), conn, shiftID)
		if err != nil {
			http.Error(w, "unable to fetch volunteers", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(volunteers)
	}
}
