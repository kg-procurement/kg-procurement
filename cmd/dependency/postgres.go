package dependency

import (
	"kg/procurement/cmd/config"
	"kg/procurement/internal/common/database"
)

func NewPostgreSQL(
	config config.PostgresConfig,
) database.DBConnector {
	return database.NewConn(
		config.Host,
		config.Username,
		config.Password,
		config.Name,
		config.Port,
	)
}
