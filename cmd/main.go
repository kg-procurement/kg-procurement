package main

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/internal/account"
	"kg/procurement/internal/product"
	"kg/procurement/internal/vendors"
	"kg/procurement/router"
	"log"
	"os"

	"github.com/benbjohnson/clock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	clock := clock.New()

	vendorSvc := vendors.NewVendorService(conn, clock)
	productSvc := product.NewProductService(conn, clock)
	accountSvc := account.NewAccountService(conn, clock)

	r := gin.Default()
	r.Use(cors.Default())
	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)
	router.NewAccountEngine(r, cfg.Routes.Account, accountSvc)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server, err: %v", err)
	}
}
