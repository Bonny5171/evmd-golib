package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	ConnectionStr string
	MaxOpenConns  int // 0 = Unlimited (default = 0)
	MaxIdleConns  int // 0 = None (default = 2)
	MaxLifetime   int // 0 = Reused forever (default = 0) - in minutes
}

func (p *PostgresDB) Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", p.ConnectionStr)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open(): %w", err)
	}

	if p.MaxOpenConns >= 0 {
		db.SetMaxOpenConns(p.MaxOpenConns)
	}

	if p.MaxIdleConns >= 0 {
		db.SetMaxIdleConns(p.MaxIdleConns)
	}

	if p.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Minute * time.Duration(p.MaxLifetime))
	}

	return db, nil
}
