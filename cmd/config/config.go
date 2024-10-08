package config

import (
	"context"
	"fmt"
	"os"
)

type Application struct {
	Common Common `mapstructure:"common" validate:"required"`
	Routes Routes `mapstructure:"routes" validate:"required"`
}

type Common struct {
	Postgres PostgresConfig `mapstructure:"postgres" validate:"required"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
}

type Routes struct {
	Vendor  VendorRoutes  `mapstructure:"vendor" validate:"required"`
	Product ProductRoutes `mapstructure:"product" validate:"required"`
}

type VendorRoutes struct {
	GetAll       string `mapstructure:"get-all" validate:"required"`
	UpdateDetail string `mapstructure:"update-detail" validate:"required"`
	GetById      string `mapstructure:"get-by-id" validate:"required"`
}

type ProductRoutes struct {
	GetProductsByVendor string `mapstructure:"get-products-by-vendor" validate:"required"`
	UpdateProduct string `mapstructure:"update-product" validate:"required"`
	UpdatePrice string `mapstructure:"update-price" validate:"required"`
}

func Load() Application {
	ctx := context.Background()
	cfgManager := NewConfigManager()

	conf := Application{}
	if err := cfgManager.Start(ctx, &conf); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config Loaded: %+v\n", conf)

	return conf
}
