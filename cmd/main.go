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

	cfg := config.Load()

	conn := dependency.NewPostgreSQL(cfg.Common.Postgres)
	defer func() {
		err := conn.Close()
		if err != nil {
			utils.Logger.Fatalf("failed to close db, err: %v", err)
		}
		_ = os.Stdout.Sync()
	}()

	awsCfg := dependency.NewAWSConfig(cfg.AWS)
	_ = mailer.NewSESProvider(*awsCfg)

	clock := clock.New()

	netSMTP := mailer.NewNativeSMTP(cfg.SMTP)

	mailerSvc := mailer.NewEmailStatusService(conn, clock)
	vendorSvc := vendors.NewVendorService(cfg, conn, clock, netSMTP, mailerSvc)
	productSvc := product.NewProductService(conn, clock)
	tokenSvc := token.NewTokenService(cfg.Token, clock)
	accountSvc := account.NewAccountService(conn, clock, tokenSvc)

	r := gin.Default()
	r.Use(cors.Default())
	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)
	router.NewAccountEngine(r, cfg.Routes.Account, accountSvc)
	router.NewEmailStatusEngine(r, cfg.Routes.EmailStatus, mailerSvc)

	if err := r.Run(":8080"); err != nil {
		utils.Logger.Fatalf("failed to run server, err: %v", err)
	}
	utils.Logger.Info("Application starts listening")
}