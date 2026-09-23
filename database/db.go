package database

import (
	"context"
	"fmt"
	"os"

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
	capacity int NOT NULL,
	start_time TIMESTAMPTZ NOT NULL,
	end_time TIMESTAMPTZ NOT NULL,
	opportunity_id int REFERENCES opportunities(id));`
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
