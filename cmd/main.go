package main

import (
	"fmt"
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
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// main is the entrypoint for the entire service
func main() {
	cfg := config.Load()

	var nrApp *newrelic.Application
	if cfg.NewRelic.Enabled {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName(cfg.NewRelic.ApplicationName),
			newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
			newrelic.ConfigDebugLogger(os.Stdout),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			fmt.Println("unable to create New Relic Application", err)
		}
		nrApp = app
		utils.ApplyNewRelicIntegration(nrApp)
	}

	conn := dependency.NewPostgreSQL(cfg.Common.Postgres)
	defer func() {
		err := conn.Close()
		if err != nil {
			utils.Logger.Fatalf("failed to close db, err: %v", err)
		}
		_ = os.Stdout.Sync()
	}()

	clock := clock.New()
	awsCfg := dependency.NewAWSConfig(cfg.AWS)

	// SMTP Providers
	_ = mailer.NewSESProvider(*awsCfg)
	_ = mailer.NewNativeSMTP(cfg.SMTP)
	gomailSMTP := mailer.NewGomailSMTP(cfg.SMTP)

	mailerSvc := mailer.NewEmailStatusService(conn, clock)
	vendorSvc := vendors.NewVendorService(cfg, conn, clock, gomailSMTP, mailerSvc)
	productSvc := product.NewProductService(conn, clock)
	tokenSvc := token.NewTokenService(cfg.Token, clock)
	accountSvc := account.NewAccountService(conn, clock, tokenSvc)

	r := gin.Default()

	r.Use(cors.Default())
	r.Use(nrgin.Middleware(nrApp))

	router.NewVendorEngine(r, cfg.Routes.Vendor, vendorSvc)
	router.NewProductEngine(r, cfg.Routes.Product, productSvc)
	router.NewAccountEngine(r, cfg.Routes.Account, accountSvc)
	router.NewEmailStatusEngine(r, cfg.Routes.EmailStatus, mailerSvc)

	if err := r.Run(":8080"); err != nil {
		utils.Logger.Fatalf("failed to run server, err: %v", err)
	}
	utils.Logger.Info("Application starts listening")
}
