package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}

func CreateSchema(conn *pgx.Conn) error {
	query := `CREATE TABLE IF NOT EXISTS volunteers (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW());`

	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS opportunities (
	id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	location TEXT NOT NULL,
	date TIMESTAMPTZ NOT NULL);`

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}
