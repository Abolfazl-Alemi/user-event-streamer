package rabbit

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user-event-streamer/pkg/logger"
)

// ConsumeMsg is responsible for consuming message from RabbitMQ
// output data in []byte format goes to Rabbit.ConsumerOutChan
func (rb *Consumer) ConsumeMsg(wg *sync.WaitGroup) {

	msg, err := rb.consumerChannel.Consume(
		rb.cfg.Connection.QueueName,
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		logger.Zap.Fatal("error in establishing consumer for rabbit", zap.Error(err), zap.String("queue_name", rb.cfg.Connection.QueueName))
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalChan
		close(rb.ConsumerOutChan)
		logger.Zap.Info("rabbit consumer out chan has been closed")
	}()

	for d := range msg {
		wg.Add(1)
		go func(data []byte) {
			defer wg.Done()
			rb.ConsumerOutChan <- data
		}(d.Body)
	}
	wg.Wait()
}

// PublishMsg is responsible for sending message to RabbitMQ
func (rb *Publisher) PublishMsg(msg []byte) error {
	err := rb.publisherChannel.Publish(rb.cfg.Connection.ExchangeName,
		rb.cfg.Connection.QueueName,
		false,
		false,
		amqp.Publishing{
			//Headers:         nil,
			ContentType: rb.cfg.Publisher.ContentType,
			//ContentEncoding: "",
			//DeliveryMode:    0,
			//Priority:        0,
			//CorrelationId:   "",
			//ReplyTo:         "",
			//Expiration:      "",
			//MessageId:       "",
			//Timestamp:       time.Time{},
			//Type:            "",
			//UserId:          "",
			//AppId:           "",
			Body: msg,
		},
	)
	if err == nil {
		return nil
	}
	return err
}
