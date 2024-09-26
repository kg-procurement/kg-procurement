package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Conn struct {
	db *sql.DB
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

func (c *Conn) Close() error {
	return c.db.Close()
}
