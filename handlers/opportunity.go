package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BosBJJ/volunteer_app/models"
)

var opportunities []models.Opportunity

func OpportunityHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		CreateOpportunity(w, req)
	case http.MethodGet:
		GetOpportunities(w, req)
	}
}

var nextOpportunityID int = 0

func CreateOpportunity(w http.ResponseWriter, req *http.Request) {
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
	opportunity.Id = nextOpportunityID
	opportunities = append(opportunities, opportunity)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(opportunity)
	nextOpportunityID++
}

func GetOpportunities(w http.ResponseWriter, req *http.Request) {
	var futureOpportunities []models.Opportunity
	now := time.Now()
	for _, opportunity := range opportunities {
		if opportunity.Time.After(now) {
			futureOpportunities = append(futureOpportunities, opportunity)
		}
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(futureOpportunities)
}
