package main

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/internal/account"
	"kg/procurement/internal/mailer"
	"kg/procurement/internal/product"
	"kg/procurement/internal/token"
	"kg/procurement/internal/vendors"
	"kg/procurement/router"
	"log"
	"os"

	"github.com/benbjohnson/clock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	u "kg/procurement/cmd/utils"
)

// main is the entrypoint for the entire service
func main() {
	logger := u.NewConsoleLogger()

	logger.Info("Starting...")
	logger.Debug("Printing stuff")
	logger.Error("Error occured")

	logger.Info("Loading configurations")
	cfg := config.Load()

	logger.Info("Creating database connection")
	conn := dependency.NewPostgreSQL(cfg.Common.Postgres)
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close db, err: %v", err)
		}
		_ = os.Stdout.Sync()
	}()

	logger.Info("Loading AWS Configurations")
	awsCfg := dependency.NewAWSConfig(cfg.AWS)
	_ = mailer.NewSESProvider(*awsCfg)

	clock := clock.New()

	logger.Info("Creating SMTP provider")
	netSMTP := mailer.NewNativeSMTP(cfg.SMTP)

	logger.Info("Preparing application services")
	vendorSvc := vendors.NewVendorService(cfg, conn, clock, netSMTP)
	productSvc := product.NewProductService(conn, clock)
	tokenSvc := token.NewTokenService(cfg.Token, clock)
	accountSvc := account.NewAccountService(conn, clock, tokenSvc)

	logger.Info("Initiating application routing engine")
	r := gin.Default()
	r.Use(cors.Default())
	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)
	router.NewAccountEngine(r, cfg.Routes.Account, accountSvc)

	logger.Info("Application starts listening")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server, err: %v", err)
	}
}
