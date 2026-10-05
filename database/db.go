package database

import (
	"context"
	"fmt"
	"os"

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

