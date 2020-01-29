package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Conns struct {
	List map[string]interface{}
}

var Connections Conns

func init() {
	Connections = Conns{}
	Connections.List = make(map[string]interface{})
}

func GetConnection(n string) (conn *sqlx.DB, err error) {
	dbq, exists := Connections.List[n]
	if !exists {
		dbq, exists = Connections.List["default"]
		if !exists {
			err = errors.New("DB queue connection not found")
			return
		}
	}

	conn = dbq.(*sqlx.DB)

	return
}

func (l *Conns) Get(n string) (conn *sqlx.DB, err error) {
	dbq, exists := l.List[n]
	if !exists {
		dbq, exists = l.List["default"]
		if !exists {
			err = errors.New("DB queue connection not found")
			return
		}
	}

	conn = dbq.(*sqlx.DB)

	return
}

func (l *Conns) Connect(n string, d DBDriver) (err error) {
	_, exists := l.List[n]
	if exists {
		return nil
	}

	db, err := d.Connect()
	if err != nil {
		return fmt.Errorf("db.Connect(): %w", err)
	}

	l.List[n] = db

	return nil
}

func (l *Conns) Close(n string) error {
	if err := l.List[n].(*sqlx.DB).Close(); err != nil {
		return fmt.Errorf("List.Close(): %w", err)
	}
	return nil
}

func (l *Conns) Check(n string) error {
	if err := l.List[n].(*sqlx.DB).Ping(); err != nil {
		return fmt.Errorf("List.Ping(): %w", err)
	}
	return nil
}
