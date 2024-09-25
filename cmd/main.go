package main

import (
	"fmt"
	"kg/procurement/cmd/config"
)

func main() {
	config := config.Load()

	fmt.Println(config.Common.Database.Postgres)
}
