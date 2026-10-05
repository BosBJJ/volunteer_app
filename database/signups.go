package database

import (
	"context"
	"fmt"

	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

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
