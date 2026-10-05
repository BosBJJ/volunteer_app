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

func ShiftHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			CreateShift(conn, w, req)
		case http.MethodGet:
			GetShifts(conn, w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func CreateShift(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	var shift models.Shift
	err := json.NewDecoder(req.Body).Decode(&shift)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if shift.OpportunityID <= 0 || shift.Capacity <= 0 || shift.Start.IsZero() || shift.End.IsZero() {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if !shift.End.After(shift.Start) {
		http.Error(w, "invalid shift end time", http.StatusBadRequest)
		return
	}
	err = database.SaveShift(req.Context(), conn, &shift)
	if err != nil {
		http.Error(w, "error saving shift to database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(&shift)
}

func GetShifts(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	shifts, err := database.ListShifts(req.Context(), conn)
	if err != nil {
		http.Error(w, "error fetching shifts", http.StatusInternalServerError)
		return
	}
	futureShifts := []models.Shift{}
	now := time.Now()
	for _, shift := range shifts {
		if shift.Start.After(now) {
			futureShifts = append(futureShifts, shift)
		}
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(futureShifts)
}

func GetShiftsByOpportunityID(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.PathValue("id")
		reqID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		shifts, err := database.ListShiftsByOpportunity(req.Context(), conn, reqID)
		if err != nil {
			http.Error(w, "error fetching shifts", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(shifts)
	}
}

func GetShiftByShiftID(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.PathValue("id")
		reqID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		shift, err := database.ListShiftByShiftId(req.Context(), conn, reqID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "invalid shift id", http.StatusNotFound)
			} else {
				http.Error(w, "error fetching shift", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(shift)
	}
}
