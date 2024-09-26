package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"user-event-streamer/pkg/logger"
)

// Publish is responsible for publishing data into kafka
func (kf *Writer) Publish(msg []byte) error {
	// if publishing data into kafka takes more than 2 seconds, ctx will be closed
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	// in case of any interruption, ctx would be notified and closed
	_, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	defer func() {
		cancel()
		stop()
	}()

	kafkaMsg := kafka.Message{
		Value: msg,
	}
	err := kf.writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		logger.Zap.Error(
			"error in publishing message into kafka",
			zap.Error(err),
			zap.String("topic", kafkaMsg.Topic),
			zap.Int("partition", kafkaMsg.Partition),
			zap.Int64("offset", kafkaMsg.Offset),
		)
		return err
	}
	return nil
}

// Consume is responsible for listening and getting new messages from kafka
// consumed messages will push into Kafka.ReaderOutChan with []byte format
func (kf *Reader) Consume(wg *sync.WaitGroup, readerRoutinesCount uint8) {
	// in case of any interruption, ctx would be notified and closed
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		stop()
		close(kf.ReaderOutChan)
	}()

	for i := uint8(0); i < readerRoutinesCount; i++ { // ReaderRoutinesCount is count of go routines for consuming messages from kafka cluster

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				m, err := kf.reader.ReadMessage(ctx)
				if err != nil {
					err2 := kf.reader.Close()
					if err2 != nil {
						logger.Zap.Error("error in closing kafka reader", zap.Error(err2), zap.String("operation", "consuming message from kafka"))
					}
					logger.Zap.Fatal(
						"error in consuming message from kafka",
						zap.Error(err),
						zap.String("topic", m.Topic),
						zap.Int("partition", m.Partition),
						zap.Int64("offset", m.Offset),
					)
				}
				kf.ReaderOutChan <- m.Value
			}
		}()
	}
	wg.Wait()
}
