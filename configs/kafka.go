package configs

type Kafka struct {
	BrokersUrl    string `yaml:"brokers_url" envconfig:"KAFKA_LEAD_VIEW_BROKERS_URL"`       // BrokersUrl example: "10.0.0.83:9092,10.0.0.84:9092". must be comma separated
	ReadTopic     string `yaml:"read_topic" envconfig:"KAFKA_LEAD_VIEW_READ_TOPIC"`         // ReadTopic could be null in case of just writing purpose
	WriteTopic    string `yaml:"write_topic" envconfig:"KAFKA_LEAD_VIEW_WRITE_TOPIC"`       // WriteTopic could be null in case of just reading purpose
	ConsumerGroup string `yaml:"consumer_group" envconfig:"KAFKA_LEAD_VIEW_CONSUMER_GROUP"` // ConsumerGroup could be null in case of just writing purpose
}
