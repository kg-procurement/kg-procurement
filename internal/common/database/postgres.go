//go:generate mockgen -typed -source=postgres.go -destination=postgres_mock.go -package=database
package database

import (
	"database/sql"
	"fmt"
	"kg/procurement/cmd/utils"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DBConnector provides simple interface for mocking purposes
type DBConnector interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	Close() error
}

type Conn struct {
	db *sqlx.DB
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

func (c *Conn) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return c.db.QueryRowx(query, args...)
}

func (c *Conn) Close() error {
	return c.db.Close()
}

// sqlx specific wrappers
func (c *Conn) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return c.db.Queryx(query, args...)
}

func (c *Conn) NamedQuery(query string, args interface{}) (*sqlx.Rows, error) {
	return c.db.NamedQuery(query, args)
}

func (c *Conn) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return c.db.NamedExec(query, arg)
}

func (c *Conn) Select(dest interface{}, query string, args ...interface{}) error {
	return c.db.Select(dest, query, args...)
}

func (c *Conn) Get(dest interface{}, query string, args ...interface{}) error {
	return c.db.Get(dest, query, args...)
}

func (c *Conn) Rebind(query string) string {
	return c.db.Rebind(query)
}

func NewConn(host, user, password, name, port string) *Conn {
	connStr := fmt.Sprintf("user=%s port=%s password=%s dbname=%s host=%s sslmode=disable",
		user, port, password, name, host)

	utils.Logger.Debug("Starts connecting db with sqlx")
	db := sqlx.MustConnect("postgres", connStr)

	if err := db.Ping(); err != nil {
		utils.Logger.Fatal(err.Error())
	}

	utils.Logger.Info("Successfully connected to database")
	log.Println("Successfully connected to database")
	return &Conn{db: db}
}
