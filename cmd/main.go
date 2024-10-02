package main

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/internal/product"
	"kg/procurement/internal/vendors"
	"kg/procurement/router"
	"log"
	"os"

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

	vendorSvc := vendors.NewVendorService(conn)
	productSvc := product.NewProductService(conn)

	r := gin.Default()

	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)

	if err := r.Run(":8000"); err != nil {
		log.Fatalf("failed to run server, err: %v", err)
	}
}
