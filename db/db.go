package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DBDriver interface {
	Connect() (*sqlx.DB, error)
}

var Conn *sqlx.DB

func Create(d DBDriver) (err error) {
	Conn, err = d.Connect()
	if err != nil {
		return fmt.Errorf("db.Connect(): %w", err)
	}
	return nil
}

// Close DB connection
func Close() error {
	if err := Conn.Close(); err != nil {
		return fmt.Errorf("Conn.Close(): %w", err)
	}
	return nil
}

// Check DB connection
func Check() error {
	if err := Conn.Ping(); err != nil {
		return fmt.Errorf("Conn.Ping(): %w", err)
	}
	return nil
}
