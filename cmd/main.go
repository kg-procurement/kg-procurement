package main

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/dependency"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/account"
	"kg/procurement/internal/mailer"
	"kg/procurement/internal/product"
	"kg/procurement/internal/token"
	"kg/procurement/internal/vendors"
	"kg/procurement/router"
	"os"

	"github.com/benbjohnson/clock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main is the entrypoint for the entire service

func main() {

	utils.Logger.Info("Starting...")

	utils.Logger.Info("Loading configurations")
	cfg := config.Load()

	utils.Logger.Info("Creating database connection")
	conn := dependency.NewPostgreSQL(cfg.Common.Postgres)
	defer func() {
		err := conn.Close()
		if err != nil {
			utils.Logger.Fatal(err.Error())
		}
		_ = os.Stdout.Sync()
	}()

	utils.Logger.Info("Loading AWS Configurations")
	awsCfg := dependency.NewAWSConfig(cfg.AWS)
	_ = mailer.NewSESProvider(*awsCfg)

	clock := clock.New()

	utils.Logger.Info("Creating SMTP provider")
	netSMTP := mailer.NewNativeSMTP(cfg.SMTP)

	utils.Logger.Info("Preparing application services")
	vendorSvc := vendors.NewVendorService(cfg, conn, clock, netSMTP)
	productSvc := product.NewProductService(conn, clock)
	tokenSvc := token.NewTokenService(cfg.Token, clock)
	accountSvc := account.NewAccountService(conn, clock, tokenSvc)

	utils.Logger.Info("Initiating application routing engines")
	r := gin.Default()
	r.Use(cors.Default())
	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)
	router.NewAccountEngine(r, cfg.Routes.Account, accountSvc)

	utils.Logger.Info("Application starts listening")
	if err := r.Run(":8080"); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}
