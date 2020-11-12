package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type connection struct {
	CanClose bool
	Conn     *sqlx.DB
	Tx       *sqlx.Tx
}

//Pool struct is a map of connections in a pool
type Pool struct {
	Connections map[string]connection
}

//NewPool returns a Pool struct initialized
func NewPool() Pool {
	p := Pool{}
	p.Connections = make(map[string]connection)
	return p
}

//Add a connection to pool
func (p *Pool) Add(name string, canClose bool, conn *sqlx.DB) (err error) {
	c := connection{
		CanClose: canClose,
		Conn:     conn,
	}

	c.Tx, err = conn.Beginx()
	if err != nil {
		return fmt.Errorf("conn.Begin(): %w", err)
	}

	p.Connections[name] = c

	return nil
}

//Conn returns connection existent in pool by name
func (p *Pool) Conn(name string) (conn *sqlx.DB) {
	c, ok := p.Connections[name]
	if ok {
		return c.Conn
	}
	return nil
}

//Tx returns a transactional connection in pool by name
func (p *Pool) Tx(name string) (conn *sqlx.Tx) {
	c, ok := p.Connections[name]
	if ok {
		return c.Tx
	}
	return nil
}

//Close closes a connection specified by name in pool
func (p *Pool) Close(name string) error {
	c, ok := p.Connections[name]
	if ok {
		if c.CanClose && c.Conn != nil {
			if err := c.Conn.Close(); err != nil {
				return fmt.Errorf("c.Conn.Close(): %w", err)
			}
		}
	}
	return nil
}

//CloseAll closes all connections existents on pool
func (p *Pool) CloseAll() error {
	for _, c := range p.Connections {
		if c.CanClose && c.Conn != nil {
			if err := c.Conn.Close(); err != nil {
				return fmt.Errorf("c.Conn.Close(): %w", err)
			}
		}
	}
	return nil
}

//Rollback rolling back a connection of pool by name
func (p *Pool) Rollback(name string) error {
	c, ok := p.Connections[name]
	if ok {
		if c.Tx != nil {
			if err := c.Tx.Rollback(); err != nil {
				return fmt.Errorf("c.Tx.Rollback(): %w", err)
			}
		}
	}
	return nil
}

//RollbackAll rolling back all connections on pool
func (p *Pool) RollbackAll() error {
	for _, c := range p.Connections {
		if c.Tx != nil {
			if err := c.Tx.Rollback(); err != nil {
				return fmt.Errorf("c.Tx.Rollback(): %w", err)
			}
		}
	}
	return nil
}

//Commit comits a set of transactions of connection on pool by name
func (p *Pool) Commit(name string) error {
	c, ok := p.Connections[name]
	if ok {
		if c.Tx != nil {
			if err := c.Tx.Commit(); err != nil {
				return fmt.Errorf("c.Tx.Commit(): %w", err)
			}
		}
	}
	return nil
}

//CommitAll commits all transactions of all connections in a pool
func (p *Pool) CommitAll() error {
	for _, c := range p.Connections {
		if c.Tx != nil {
			if err := c.Tx.Commit(); err != nil {
				return fmt.Errorf("c.Tx.Commit(): %w", err)
			}
		}
	}
	return nil
}

//GetConnDB creates and returns a new connected connection
func GetConnDB(dsn string) (*sqlx.DB, error) {
	var (
		dbMaxOpenConns int = 5
		dbMaxIdleConns int = 1
		dbMaxLifeTime  int = 10
	)

	if len(os.Getenv("GOWORKER_DB_MAXOPENCONNS")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GOWORKER_DB_MAXOPENCONNS")); e != nil {
			dbMaxOpenConns = v
		}
	}

	if len(os.Getenv("GOWORKER_DB_MAXIDLECONNS")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GOWORKER_DB_MAXIDLECONNS")); e != nil {
			dbMaxIdleConns = v
		}
	}

	if len(os.Getenv("GOWORKER_DB_MAXLIFETIME")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GOWORKER_DB_MAXLIFETIME")); e != nil {
			dbMaxLifeTime = v
		}
	}

	driver := &PostgresDB{
		ConnectionStr: dsn,
		MaxOpenConns:  dbMaxOpenConns,
		MaxIdleConns:  dbMaxIdleConns,
		MaxLifetime:   dbMaxLifeTime,
	}

	db, err := driver.Connect()
	if err != nil {
		return nil, fmt.Errorf("db.Connect(): %w", err)
	}

	return db, nil
}
