package config

import (
	"context"
	"fmt"
	"os"
)

type Application struct {
	Common Common `mapstructure:"common" validate:"required"`
	Routes Routes `mapstructure:"routes" validate:"required"`
	Token  Token  `mapstructure:"token" validate:"required"`
	SMTP   SMTP   `mapstructure:"smtp" validate:"required"`
	AWS    AWS    `mapstructure:"aws" validate:"required"`
}

type SMTP struct {
	Host         string `mapstructure:"host" validate:"required"`
	Port         string `mapstructure:"port" validate:"required"`
	SenderName   string `mapstructure:"sender_name" validate:"required"`
	AuthEmail    string `mapstructure:"auth_email" validate:"required"`
	AuthPassword string `mapstructure:"auth_password" validate:"required"`
}

type AWS struct {
	AccessKey       string `mapstructure:"access-key" validate:"required"`
	SecretAccessKey string `mapstructure:"secret-access-key" validate:"required"`
	Region          string `mapstructure:"region" validate:"required"`
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
	Account AccountRoutes `mapstructure:"account" validate:"required"`
}

type VendorRoutes struct {
	GetAll              string `mapstructure:"get-all" validate:"required"`
	UpdateDetail        string `mapstructure:"update-detail" validate:"required"`
	GetById             string `mapstructure:"get-by-id" validate:"required"`
	GetLocations        string `mapstructure:"get-locations" validate:"required"`
	EmailBlast          string `mapstructure:"email-blast" validate:"required"`
	AutomatedEmailBlast string `mapstructure:"automated-email-blast" validate:"required"`
}

type ProductRoutes struct {
	GetProductsByVendor string `mapstructure:"get-products-by-vendor" validate:"required"`
	UpdateProduct       string `mapstructure:"update-product" validate:"required"`
	UpdatePrice         string `mapstructure:"update-price" validate:"required"`
}

type Token struct {
	Secret string `mapstructure:"secret" validate:"required"`
}

type AccountRoutes struct {
	Register string `mapstructure:"register" validate:"required"`
	Login    string `mapstructure:"login" validate:"required"`
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
