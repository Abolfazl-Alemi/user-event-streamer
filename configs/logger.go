package configs

// Logger should better get just two values (production | development)
type Logger struct {
	// 1. production -> INFO
	// 2. development -> DEBUG
	// else -> INFO
	Level string `yaml:"level" envconfig:"ZAP_LOGGER_LEVEL"`
}
