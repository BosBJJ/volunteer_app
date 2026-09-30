package models

import "time"

type Attendance struct {
	Id          int        `json:"id"`
	ShiftID     int        `json:"shift_id"`
	VolunteerID int        `json:"volunteer_id"`
	CheckIn     time.Time  `json:"check_in"`
	CheckOut    *time.Time `json:"check_out"`
}
