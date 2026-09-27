package models

type Signup struct {
	Id          int `json:"id"`
	VolunteerID int `json:"volunteer_id"`
	ShiftID     int `json:"shift_id"`
}
