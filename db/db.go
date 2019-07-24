package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBDriver interface {
	Connect() (*sqlx.DB, error)
}

var Conn *sqlx.DB

func Create(d DBDriver) (err error) {
	Conn, err = d.Connect()
	if err != nil {
		return errors.Wrap(err, "db.Connect()")
	}
	return nil
}

// Close DB connection
func Close() error {
	if err := Conn.Close(); err != nil {
		return errors.Wrap(err, "Conn.Close()")
	}
	return nil
}

// Check DB connection
func Check() error {
	if err := Conn.Ping(); err != nil {
		return errors.Wrap(err, "Conn.Ping()")
	}
	return nil
}
