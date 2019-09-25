package model

type Stack struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	DSN  string `db:"dsn"`
}
