package configs

var cfg Config

// Config includes all configs structs
type Config struct {
	Redis  Redis  `yaml:"redis"`
	Kafka  Kafka  `yaml:"kafka"`
	Logger Logger `yaml:"logger"`
	Server Server `yaml:"server"`
	Rabbit Rabbit `yaml:"rabbit"`
}
