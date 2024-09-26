package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"strings"
	"user-event-streamer/pkg/logger"
)

const (
	WriterMaxAttempts int   = 10
	WriterBatchSize   int   = 100
	WriterBatchBytes  int64 = 1048576
)

type Writer struct {
	writer *kafka.Writer
}

type Reader struct {
	reader        *kafka.Reader
	ReaderOutChan chan []byte // ReaderOutChan is for sending consumed messages from kafka
}

func NewKafkaConsumer(brokersAddress string, readTopic string, consumerGroup string) *Reader {
	return &Reader{
		reader:        NewKafkaReader(brokersAddress, readTopic, consumerGroup),
		ReaderOutChan: make(chan []byte),
	}
}

func NewKafkaPublisher(brokersAddress string, writeTopic string) *Writer {
	return &Writer{
		writer: NewKafkaWriter(brokersAddress, writeTopic),
	}
}

func NewKafkaWriter(brokersAddress string, topic string) *kafka.Writer {
	addr := strings.Split(brokersAddress, ",")
	w := &kafka.Writer{
		Addr:        kafka.TCP(addr...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		MaxAttempts: WriterMaxAttempts,
		//WriteBackoffMin:        0,
		//WriteBackoffMax:        0,
		BatchSize:  WriterBatchSize,
		BatchBytes: WriterBatchBytes,
		//BatchTimeout:           0,
		//ReadTimeout:            0,
		//WriteTimeout:           0,
		RequiredAcks: kafka.RequireNone,
		//RequireNone (0)  fire-and-forget, do not wait for acknowledgements;
		//RequireOne  (1)  wait for the leader to acknowledge the writes
		//RequireAll  (-1) wait for the full ISR to acknowledge the writes
		Async: false,
		//Completion:             nil,
		Compression: kafka.Snappy,
		//Logger:      kafka.LoggerFunc(kafkaLogger),
		ErrorLogger: kafka.LoggerFunc(kafkaErrorLogger),
		//Transport:              nil,
		AllowAutoTopicCreation: true,
	}
	return w
}

func NewKafkaReader(brokersAddress string, topic string, consumerGroup string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokersAddress, ","),
		GroupID: consumerGroup,
		//GroupTopics:            nil,
		Topic: topic,
		//Dialer:                 nil,
		//QueueCapacity:          0,
		//MinBytes:               0,
		//MaxBytes:               0,
		//MaxWait:                0,
		//ReadBatchTimeout:       0,
		//ReadLagInterval:        0,
		//GroupBalancers:         nil,
		//HeartbeatInterval:      0,
		//CommitInterval: 0,
		//PartitionWatchInterval: 0,
		//WatchPartitionChanges:  false,
		//SessionTimeout:         0,
		//RebalanceTimeout:       0,
		//JoinGroupBackoff:       0,
		//RetentionTime:          0,
		StartOffset: kafka.LastOffset,
		//ReadBackoffMin:         0,
		//ReadBackoffMax:         0,
		//Logger:      kafka.LoggerFunc(kafkaLogger),
		ErrorLogger: kafka.LoggerFunc(kafkaErrorLogger),
		//IsolationLevel:         0,
		//MaxAttempts:            0,
		//OffsetOutOfRangeError:  false,
	})
	return r
}

func kafkaErrorLogger(msg string, a ...interface{}) {
	logger.Zap.Error("kafka custom error logger", zap.String("error", fmt.Sprintf("%v", a)))
}
