//go:generate mockgen -typed -source=postgres.go -destination=postgres_mock.go -package=database
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DBConnector provides simple interface for mocking purposes
type DBConnector interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

type Conn struct {
	db *sql.DB
}

func (c *Conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

func (c *Conn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c *Conn) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.db.QueryRow(query, args...)
}

func (c *Conn) Close() error {
	return c.db.Close()
}

func NewConn(host, user, password, name, port string) *Conn {
	connStr := fmt.Sprintf("user=%s port=%s password=%s dbname=%s host=%s sslmode=disable",
		user, port, password, name, host)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	log.Println("Successfully connected to database")
	return &Conn{db: db}
}
