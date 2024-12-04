package dependency

import (
	"kg/procurement/cmd/config"
	"kg/procurement/internal/common/database"
)

func NewRedisCache(config config.RedisClientConfig) database.RedisClientInterface {
	return database.NewRedisClient(config.Address, config.Password)
}
