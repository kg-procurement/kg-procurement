package main

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/internal/account"
	"kg/procurement/internal/product"
	"kg/procurement/internal/smtp_provider"
	"kg/procurement/internal/token"
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

	awsCfg := dependency.NewAWSConfig(cfg.AWS)
	_ = smtp_provider.NewSESProvider(*awsCfg)

	clock := clock.New()

	netSMTP := smtp_provider.NewNativeSMTP(cfg.SMTP)

	vendorSvc := vendors.NewVendorService(cfg, conn, clock, netSMTP)
	productSvc := product.NewProductService(conn, clock)
	tokenSvc := token.NewTokenService(cfg.Token, clock)
	accountSvc := account.NewAccountService(conn, clock, tokenSvc)

	r := gin.Default()
	r.Use(cors.Default())
	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)
	router.NewAccountEngine(r, cfg.Routes.Account, accountSvc)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("failed to run server, err: %v", err)
	}
}
