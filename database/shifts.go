package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

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
	query := `SELECT id, capacity, start_time, end_time, opportunity_id FROM shifts`
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
	query := `SELECT id, capacity, start_time, end_time, opportunity_id FROM shifts WHERE id = $1`
	row := conn.QueryRow(ctx, query, shiftId)
	err := row.Scan(&shift.Id, &shift.Capacity, &shift.Start, &shift.End, &shift.OpportunityID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Shift{}, fmt.Errorf("no shift with id %d: %w", shiftId, err)
		} else {
			return models.Shift{}, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return shift, nil
}

func ListShiftsByOpportunity(ctx context.Context, conn *pgx.Conn, opportunityID int) ([]models.Shift, error) {
	shifts := []models.Shift{}
	query := `SELECT id, capacity, start_time, end_time, opportunity_id FROM shifts WHERE opportunity_id = $1`
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
