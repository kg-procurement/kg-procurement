package main

import (
	"fmt"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"log"
	"os"
)

// main is the entrypoint for the entire service
func main() {
	cfg := config.Load()

	conn := dependency.NewPostgreSQL(cfg.Common.Postgres)
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close db, err: %v", err)
		}
		_ = os.Stdout.Sync()
	}()

	fmt.Println("Welcomee")
}
