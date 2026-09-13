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

func OpportunityHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			CreateOpportunity(conn, w, req)
		case http.MethodGet:
			GetOpportunities(conn, w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func CreateOpportunity(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	var opportunity models.Opportunity
	err := json.NewDecoder(req.Body).Decode(&opportunity)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if opportunity.Title == "" || opportunity.Location == "" || opportunity.Description == "" {
		http.Error(w, "Title, Location and Description are required.", http.StatusBadRequest)
		return
	}
	err = database.SaveOpportunity(conn, &opportunity)
	if err != nil {
		http.Error(w, "error saving opportunity to database", http.StatusInternalServerError)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(opportunity)
}

func GetOpportunities(conn *pgx.Conn, w http.ResponseWriter, req *http.Request) {
	opportunities, err := database.ListOpportunities(conn)
	if err != nil {
		http.Error(w, "error fetching opportunities", http.StatusInternalServerError)
		return
	}
	futureOpportunities := []models.Opportunity{}
	now := time.Now()
	for _, opportunity := range opportunities {
		if opportunity.Date.After(now) {
			futureOpportunities = append(futureOpportunities, opportunity)
		}
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(futureOpportunities)
}

func GetOpportunityByID(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.PathValue("id")
		reqID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		opportunity, err := database.ListOpportunityByID(conn, reqID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "invalid opportunity id", http.StatusNotFound)
			} else {
				http.Error(w, "error fetching opportunity", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(opportunity)
	}
}
