package configs

// Redis keeps all the configs for connecting to redis sentinel servers
type Redis struct {
	SentinelAddress string `yaml:"sentinel_address" envconfig:"SENTINEL_ADDRESS"` // SentinelAddress is comma separated (sentinel addresses)
	Db              int    `yaml:"db" envconfig:"REDIS_DB"`                       // Db is the number of redis database want to use
	MasterName      string `yaml:"master_name" envconfig:"REDIS_MASTER_NAME"`     // MasterName is the redis master name which can be found in redis-cli info.
}
