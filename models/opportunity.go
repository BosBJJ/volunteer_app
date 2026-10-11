package models

import "time"

type Opportunity struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Date        time.Time `json:"date"`
}

type OpportunityUpdate struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}
