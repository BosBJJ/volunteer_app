package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

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
	query := `SELECT id, name, email, registered_at FROM volunteers`
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
	query := `SELECT id, name, email, registered_at FROM volunteers WHERE id = $1`
	row := conn.QueryRow(ctx, query, reqID)
	err := row.Scan(&volunteer.Id, &volunteer.Name, &volunteer.Email, &volunteer.RegisteredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
