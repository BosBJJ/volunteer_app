package models

import "time"


type Shift struct {
	OpportunityID int       `json:"opportunity_id"`
	Id            int       `json:"id"`
	Capacity      int       `json:"capacity"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
}
