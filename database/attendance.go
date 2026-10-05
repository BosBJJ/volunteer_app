package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

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
