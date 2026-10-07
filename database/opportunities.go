package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/BosBJJ/volunteer_app/models"
	"github.com/jackc/pgx/v5"
)

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
	query := `SELECT id, title, description, location, date FROM opportunities`
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
	query := `SELECT id, title, description, location, date FROM opportunities WHERE id = $1`
	row := conn.QueryRow(ctx, query, reqID)
	err := row.Scan(&opportunity.Id, &opportunity.Title, &opportunity.Description, &opportunity.Location, &opportunity.Date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Opportunity{}, fmt.Errorf("no opportunity with id %d: %w", reqID, err)
		} else {
			return models.Opportunity{}, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return opportunity, nil
}

var ErrOpportunityNotFound = errors.New("opportunity not found")

func DeleteOpportunity(ctx context.Context, conn *pgx.Conn, opportunityID int) error {
	query := `DELETE FROM opportunities WHERE id = $1`
	tag, err := conn.Exec(ctx, query, opportunityID)
	if err != nil {
		return fmt.Errorf("unable to delete opportunity: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrOpportunityNotFound
	}
	return nil
}
