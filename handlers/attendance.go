package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/BosBJJ/volunteer_app/database"
	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

func ShiftCheckIn(conn *pgx.Conn) http.HandlerFunc {
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
		exists, err := database.SignupExists(req.Context(), conn, volunteerID, shiftID)
		if err != nil {
			http.Error(w, "unable to fetch signups", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "user not registered for specified shift", http.StatusConflict)
			return
		}
		checkInTime := time.Now()
		attendance := models.Attendance{
			ShiftID:     shiftID,
			VolunteerID: volunteerID,
			CheckIn:     checkInTime,
		}
		err = database.SaveAttendance(req.Context(), conn, &attendance)
		if err != nil {
			http.Error(w, "unable to save attendance", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(attendance)
	}
}

func ShiftCheckOut(conn *pgx.Conn) http.HandlerFunc {
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
		checkOutTime := time.Now()
		err = database.UpdateAttendance(req.Context(), conn, checkOutTime, volunteerID, shiftID)
		if err != nil {
			if errors.Is(err, database.ErrAttendanceNotFound) {
				http.Error(w, "attendance record not found", http.StatusConflict)
				return
			}
			http.Error(w, "unable to update attendance", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode("Complete")
	}
}
