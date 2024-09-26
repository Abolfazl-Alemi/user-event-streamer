package configs

// Rabbit keeps all rabbit client configs:
type Rabbit struct {
	Server     RabbitServer     `yaml:"server"`      // Server: includes connection with rabbit server info -> host, port, user, password
	Connection RabbitConnection `yaml:"connection"`  // Connection: includes exchange and queue info
	DeadLetter DeadLetter       `yaml:"dead_letter"` // DeadLetter: in case of having dead letter for queues, this class keeps the data.
	Publisher  Publisher        `yaml:"publisher"`
}

// RabbitServer includes connection with rabbit server info -> host, port, user, password
type RabbitServer struct {
	Host     string `yaml:"host" envconfig:"RABBIT_HOST"`
	Port     uint   `yaml:"port" envconfig:"RABBIT_PORT"`
	User     string `yaml:"user" envconfig:"RABBIT_USER"`
	Password string `yaml:"pass" envconfig:"RABBIT_PASSWORD"`
}

// RabbitConnection includes exchange and queue info
type RabbitConnection struct {
	ExchangeName     string `yaml:"exchange_name" envconfig:"RABBIT_EXCHANGE_NAME"`
	ExchangeType     string `yaml:"exchange_type" envconfig:"RABBIT_EXCHANGE_TYPE"` // ExchangeType: valid values -> direct, topic, fanout, headers
	RoutingKey       string `yaml:"routing_key" envconfig:"RABBIT_ROUTING_KEY"`
	QueueName        string `yaml:"queue_name" envconfig:"RABBIT_QUEUE_NAME"`
	DelayToReconnect int    `yaml:"reconnect_delay_second" envconfig:"RABBIT_DELAY_TO_RECONNECT"` // DelayToReconnect: in case of retries for connect to the server, this is the delay in seconds
	PrefetchCount    int    `yaml:"prefetch_count" envconfig:"RABBIT_PREFETCH_COUNT"`
}

// DeadLetter  in case of having dead letter for queues, this class keeps the data.
type DeadLetter struct {
	ExchangeName string `yaml:"exchange_name" envconfig:"RABBIT_DEAD_LETTER_EXCHANGE_NAME"`
	RoutingKey   string `yaml:"routing_key" envconfig:"RABBIT_DEAD_LETTER_ROUTING_KEY"`
	Ttl          int    `yaml:"ttl" envconfig:"RABBIT_DEAD_LETTER_TTL"` // Ttl is in seconds
}

type Publisher struct {
	ContentType string `yaml:"content_type"`
}
