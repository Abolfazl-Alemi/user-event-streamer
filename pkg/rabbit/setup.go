package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"user-event-streamer/configs"
	"user-event-streamer/pkg/logger"
)

type Consumer struct {
	cfg             *configs.Rabbit
	consumerChannel *amqp.Channel // consumerChannel will declare Exchange and Queue and bind the RoutingKey to the Queue
	ConsumerOutChan chan []byte   // ConsumerOutChan is responsible for sending consumed messages from kafka to Wrapper.RabbitConsumerToElasticGetId
}

type Publisher struct {
	cfg              *configs.Rabbit
	publisherChannel *amqp.Channel // publisherChannel just bind the RoutingKey to the Queue
}

func NewRabbitConsumer(cfg *configs.Rabbit) *Consumer {
	con := NewConnection(cfg)
	return &Consumer{
		cfg:             cfg,
		consumerChannel: connectToChannel(con, cfg, true),
		ConsumerOutChan: make(chan []byte),
	}
}

func NewRabbitPublisher(cfg *configs.Rabbit) *Publisher {
	con := NewConnection(cfg)
	return &Publisher{
		cfg:              cfg,
		publisherChannel: connectToChannel(con, cfg, false),
	}
}

func connectToChannel(rb *amqp.Connection, cfg *configs.Rabbit, isConsumer bool) *amqp.Channel {
	ch, err := rb.Channel()
	if err != nil {
		logger.Zap.Fatal(err.Error())
	}

	err = ch.Confirm(false)
	if err != nil {
		logger.Zap.Fatal(err.Error())
	}

	if isConsumer == true { // just for consumers we need to declare Exchange and Queue
		if err := ch.ExchangeDeclare(
			cfg.Connection.ExchangeName,
			cfg.Connection.ExchangeType,
			false, // Durable
			false, // AutoDelete
			false, // Internal
			false, // NoWait
			nil,   // Arguments
		); err != nil {
			logger.Zap.Fatal(
				"error in declaring the Rabbit exchange",
				zap.Error(err),
				zap.String("exchange_name", cfg.Connection.ExchangeName),
			)
		}

		_, err = ch.QueueDeclare(
			cfg.Connection.QueueName,
			true,  // Durable
			false, // Delete when unused
			false, // Exclusive
			false, // No-wait
			amqp.Table{ // Arguments
				"x-dead-letter-exchange":    cfg.DeadLetter.ExchangeName, // Specify the Dead Letter Exchange
				"x-dead-letter-routing-key": cfg.DeadLetter.RoutingKey,
				"x-message-ttl":             cfg.DeadLetter.Ttl * 1000,
			},
		)
		if err != nil {
			logger.Zap.Fatal(
				"error in declaring the Rabbit queue",
				zap.Error(err),
				zap.String("exchange_name", cfg.Connection.ExchangeName),
				zap.String("queue_name", cfg.Connection.QueueName),
			)
		}
	}

	if err := ch.QueueBind(
		cfg.Connection.QueueName,
		cfg.Connection.RoutingKey,
		cfg.Connection.ExchangeName,
		false,
		nil,
	); err != nil {
		logger.Zap.Fatal(
			"error in binding routing key to the queue",
			zap.Error(err),
			zap.String("exchange_name", cfg.Connection.ExchangeName),
			zap.String("queue_name", cfg.Connection.QueueName),
			zap.String("routing_key", cfg.Connection.RoutingKey),
		)
	}
	// Qos controls how many messages or how many bytes the server will try to keep on the network for consumers before receiving delivery acks.
	err = ch.Qos(cfg.Connection.PrefetchCount, 0, false)
	if err != nil {
		logger.Zap.Fatal(
			"error in Rabbit Qos",
			zap.Error(err),
			zap.String("exchange_name", cfg.Connection.ExchangeName),
			zap.String("queue_name", cfg.Connection.QueueName),
			zap.String("routing_key", cfg.Connection.RoutingKey),
		)
	}

	return ch
}

func NewConnection(cfg *configs.Rabbit) *amqp.Connection {

	for i := 1; i <= 3; i++ {
		logger.Zap.Info("Connecting to Rabbit", zap.Int("iteration", i))
		hostURL := fmt.Sprintf("amqp://%v:%v@%v:%v", cfg.Server.User, cfg.Server.Password, cfg.Server.Host, cfg.Server.Port)
		conn, err := amqp.Dial(hostURL)

		if err == nil {
			logger.Zap.Info("Connected to Rabbit", zap.String("rabbit_addr", hostURL))
			return conn
		}
		logger.Zap.Error("error in connecting to rabbit", zap.Error(err), zap.String("rabbit_addr", hostURL))
	}
	logger.Zap.Fatal("error in connecting to rabbit after all retries")
	return nil
}
