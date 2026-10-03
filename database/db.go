package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}

func CreateSchema(ctx context.Context, conn *pgx.Conn) error {
	query := `CREATE TABLE IF NOT EXISTS volunteers (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW());`
	_, err := conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS opportunities (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	location TEXT NOT NULL,
	date TIMESTAMPTZ NOT NULL);`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS shifts (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	capacity INTEGER NOT NULL,
	start_time TIMESTAMPTZ NOT NULL,
	end_time TIMESTAMPTZ NOT NULL,
	opportunity_id INTEGER NOT NULL REFERENCES opportunities(id));`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS shift_signups (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	volunteer_id INTEGER NOT NULL REFERENCES volunteers(id),
	shift_id INTEGER NOT NULL REFERENCES shifts(id),
	UNIQUE (volunteer_id, shift_id));`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS attendance (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	shift_id INTEGER NOT NULL REFERENCES shifts(id),
	volunteer_id INTEGER NOT NULL REFERENCES volunteers(id),
	check_in TIMESTAMPTZ NOT NULL,
	check_out TIMESTAMPTZ,
	UNIQUE (volunteer_id, shift_id));`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func SaveVolunteer(ctx context.Context, conn *pgx.Conn, volunteer *models.Volunteer) error {
	query := `INSERT INTO volunteers (name, email) VALUES ($1, $2) RETURNING id, registered_at`
	row := conn.QueryRow(ctx, query, volunteer.Name, volunteer.Email)
	err := row.Scan(&volunteer.Id, &volunteer.RegisteredAt)
	if err != nil {
		return fmt.Errorf("unable to save volunteer to database: %w", err)
	}
	return nil
}

func ListVolunteers(ctx context.Context, conn *pgx.Conn) ([]models.Volunteer, error) {
	volunteers := []models.Volunteer{}
	query := `SELECT * FROM volunteers`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching volunteers: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var volunteer models.Volunteer
		err = rows.Scan(&volunteer.Id, &volunteer.Name, &volunteer.Email, &volunteer.RegisteredAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning volunteer: %w", err)
		}
		volunteers = append(volunteers, volunteer)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading volunteers: %w", err)
	}
	return volunteers, nil
}

func ListVolunteerByID(ctx context.Context, conn *pgx.Conn, reqID int) (models.Volunteer, error) {
	volunteer := models.Volunteer{}
	query := `SELECT * FROM volunteers WHERE id = $1`
	row := conn.QueryRow(ctx, query, reqID)
	err := row.Scan(&volunteer.Id, &volunteer.Name, &volunteer.Email, &volunteer.RegisteredAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Volunteer{}, fmt.Errorf("no volunteer with id %d: %w", reqID, err)
		} else {
			return models.Volunteer{}, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return volunteer, nil
}

func VolunteerExists(ctx context.Context, conn *pgx.Conn, volunteerID int) (bool, error) {
	var res bool
	query := `SELECT EXISTS(SELECT 1 FROM volunteers WHERE id = $1);`
	row := conn.QueryRow(ctx, query, volunteerID)
	err := row.Scan(&res)
	if err != nil {
		return res, fmt.Errorf("unable to scan for volunteer: %w", err)
	}
	return res, nil
}

func SaveOpportunity(ctx context.Context, conn *pgx.Conn, opportunity *models.Opportunity) error {
	query := `INSERT INTO opportunities (title, description, location, date) VALUES ($1, $2, $3, $4) RETURNING id`
	row := conn.QueryRow(ctx, query, opportunity.Title, opportunity.Description, opportunity.Location, opportunity.Date)
	err := row.Scan(&opportunity.Id)
	if err != nil {
		return fmt.Errorf("unable to save opportunity to database: %w", err)
	}
	return nil
}

func ListOpportunities(ctx context.Context, conn *pgx.Conn) ([]models.Opportunity, error) {
	query := `SELECT * FROM opportunities`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching opportunities: %w", err)
	}
	defer rows.Close()
	opportunities := []models.Opportunity{}
	for rows.Next() {
		var opportunity models.Opportunity
		err := rows.Scan(&opportunity.Id, &opportunity.Title, &opportunity.Description, &opportunity.Location, &opportunity.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning opportunity: %w", err)
		}
		opportunities = append(opportunities, opportunity)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading opportunities: %w", err)
	}
	return opportunities, nil
}

func ListOpportunityByID(ctx context.Context, conn *pgx.Conn, reqID int) (models.Opportunity, error) {
	opportunity := models.Opportunity{}
	query := `SELECT * FROM opportunities WHERE id = $1`
	row := conn.QueryRow(ctx, query, reqID)
	err := row.Scan(&opportunity.Id, &opportunity.Title, &opportunity.Description, &opportunity.Location, &opportunity.Date)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Opportunity{}, fmt.Errorf("no opportunity with id %d: %w", reqID, err)
		} else {
			return models.Opportunity{}, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return opportunity, nil
}

func SaveShift(ctx context.Context, conn *pgx.Conn, shift *models.Shift) error {
	query := `INSERT INTO shifts (opportunity_id, capacity, start_time, end_time) VALUES ($1, $2, $3, $4) returning id`
	row := conn.QueryRow(ctx, query, shift.OpportunityID, shift.Capacity, shift.Start, shift.End)
	err := row.Scan(&shift.Id)
	if err != nil {
		return fmt.Errorf("unable to save shift to database: %w", err)
	}
	return nil
}

func ListShifts(ctx context.Context, conn *pgx.Conn) ([]models.Shift, error) {
	query := `SELECT * FROM shifts`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching shifts: %w", err)
	}
	defer rows.Close()
	shifts := []models.Shift{}
	for rows.Next() {
		var shift models.Shift
		err := rows.Scan(&shift.Id, &shift.Capacity, &shift.Start, &shift.End, &shift.OpportunityID)
		if err != nil {
			return nil, fmt.Errorf("error scanning shift: %w", err)
		}
		shifts = append(shifts, shift)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading shifts: %w", err)
	}
	return shifts, nil
}

func ListShiftByShiftId(ctx context.Context, conn *pgx.Conn, shiftId int) (models.Shift, error) {
	shift := models.Shift{}
	query := `SELECT id, opportunity_id, capacity, start_time, end_time FROM shifts WHERE id = $1`
	row := conn.QueryRow(ctx, query, shiftId)
	err := row.Scan(&shift.Id, &shift.OpportunityID, &shift.Capacity, &shift.Start, &shift.End)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Shift{}, fmt.Errorf("no shift with id %d: %w", shiftId, err)
		} else {
			return models.Shift{}, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return shift, nil
}

func ListShiftsByOpportunity(ctx context.Context, conn *pgx.Conn, opportunityID int) ([]models.Shift, error) {
	shifts := []models.Shift{}
	query := `SELECT * FROM shifts WHERE opportunity_id = $1`
	rows, err := conn.Query(ctx, query, opportunityID)
	if err != nil {
		return nil, fmt.Errorf("error fetching shifts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var shift models.Shift
		err := rows.Scan(&shift.Id, &shift.Capacity, &shift.Start, &shift.End, &shift.OpportunityID)
		if err != nil {
			return nil, fmt.Errorf("error scanning shift: %w", err)
		}
		shifts = append(shifts, shift)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading shifts: %w", err)
	}
	return shifts, nil
}

func SaveSignup(ctx context.Context, conn *pgx.Conn, signup *models.Signup) error {
	query := `INSERT INTO shift_signups (volunteer_id, shift_id) VALUES ($1, $2) RETURNING id`
	row := conn.QueryRow(ctx, query, signup.VolunteerID, signup.ShiftID)
	err := row.Scan(&signup.Id)
	if err != nil {
		return fmt.Errorf("unable to save signup to database: %w", err)
	}
	return nil
}

func SignupExists(ctx context.Context, conn *pgx.Conn, volunteerID, shiftID int) (bool, error) {
	var res bool
	query := `SELECT EXISTS(SELECT 1 FROM shift_signups WHERE volunteer_id = $1 AND shift_id = $2);`
	row := conn.QueryRow(ctx, query, volunteerID, shiftID)
	err := row.Scan(&res)
	if err != nil {
		return res, fmt.Errorf("unable to scan for signup: %w", err)
	}
	return res, nil
}

func CountSignups(ctx context.Context, conn *pgx.Conn, shiftID int) (int, error) {
	var res int
	query := `SELECT COUNT(*) FROM shift_signups WHERE shift_id = $1`
	row := conn.QueryRow(ctx, query, shiftID)
	err := row.Scan(&res)
	if err != nil {
		return 0, fmt.Errorf("unable to fetch signup count: %w", err)
	}
	return res, nil
}

func ListSignupsByShift(ctx context.Context, conn *pgx.Conn, shiftId int) ([]models.Volunteer, error) {
	volunteers := []models.Volunteer{}
	query := `SELECT v.id, v.name, v.email, v.registered_at FROM volunteers v JOIN shift_signups s ON v.id = s.volunteer_id WHERE s.shift_id = $1`
	rows, err := conn.Query(ctx, query, shiftId)
	if err != nil {
		return nil, fmt.Errorf("error fetching signups: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var volunteer models.Volunteer
		err = rows.Scan(&volunteer.Id, &volunteer.Name, &volunteer.Email, &volunteer.RegisteredAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning signup: %w", err)
		}
		volunteers = append(volunteers, volunteer)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading signups: %w", err)
	}
	return volunteers, nil
}

func SaveAttendance(ctx context.Context, conn *pgx.Conn, attendance *models.Attendance) error {
	query := `INSERT INTO attendance (shift_id, volunteer_id, check_in) VALUES ($1, $2, $3) RETURNING id`
	row := conn.QueryRow(ctx, query, attendance.ShiftID, attendance.VolunteerID, attendance.CheckIn)
	err := row.Scan(&attendance.Id)
	if err != nil {
		return fmt.Errorf("unable to add attendance to database: %w", err)
	}
	return nil
}

var ErrAttendanceNotFound = errors.New("no attendance record found")

func UpdateAttendance(ctx context.Context, conn *pgx.Conn, endTime time.Time, volunteerID, shiftID int) error {
	query := `UPDATE attendance SET check_out = $1 WHERE volunteer_id = $2 AND shift_id = $3 AND check_out IS NULL`
	tag, err := conn.Exec(ctx, query, endTime, volunteerID, shiftID)
	if err != nil {
		return fmt.Errorf("unable to update attendance in database: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrAttendanceNotFound
	}
	return nil
}

func GetAttendanceByVolunteerID(ctx context.Context, conn *pgx.Conn, volunteerID int) ([]models.Attendance, error) {
	attendanceRecords := []models.Attendance{}
	query := `SELECT id, shift_id, volunteer_id, check_in, check_out FROM attendance WHERE volunteer_id = $1`
	rows, err := conn.Query(ctx, query, volunteerID)
	if err != nil {
		return nil, fmt.Errorf("error fetching attendance records: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var currentAttendance models.Attendance
		err = rows.Scan(&currentAttendance.Id, &currentAttendance.ShiftID, &currentAttendance.VolunteerID, &currentAttendance.CheckIn, &currentAttendance.CheckOut)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance record: %w", err)
		}
		attendanceRecords = append(attendanceRecords, currentAttendance)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading attendance records: %w", err)
	}
	return attendanceRecords, nil
}

func GetTotalHoursByVolunteerID(ctx context.Context, conn *pgx.Conn, volunteerID int) (float64, error) {
	var hours float64
	query := `SELECT COALESCE(SUM(EXTRACT(EPOCH FROM(check_out - check_in)))/ 3600, 0)
	FROM attendance WHERE volunteer_id = $1 AND check_out IS NOT NULL;`
	row := conn.QueryRow(ctx, query, volunteerID)
	err := row.Scan(&hours)
	if err != nil {
		return 0, fmt.Errorf("error scanning total hours: %w", err)
	}
	return hours, nil
}
