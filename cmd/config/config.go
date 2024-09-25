package config

import (
	"context"
	"fmt"
	"os"
)

type Application struct {
	Common Common `mapstructure:"common" validate:"required"`
}

type Common struct {
	Database Database `mapstructure:"database" validate:"required"`
}

type Database struct {
	Postgres PostgresConfig `mapstructure:"postgres" validate:"required"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
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
