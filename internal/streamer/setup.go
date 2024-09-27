package streamer

import (
	"user-event-streamer/pkg/kafka"
	"user-event-streamer/pkg/monitoring"
	"user-event-streamer/pkg/rabbit"
	"user-event-streamer/pkg/udp"
)

type Streamer struct {
	udp              *udp.Server
	kf               *kafka.Writer
	pr               *monitoring.ProMetrics
	exceptionHandler *ExceptionHandler
}

func NewStreamer(udp *udp.Server, kf *kafka.Writer, pr *monitoring.ProMetrics, rbc *rabbit.Consumer, rbp *rabbit.Publisher, kfExc *kafka.Writer) *Streamer {
	return &Streamer{
		udp:              udp,
		kf:               kf,
		pr:               pr,
		exceptionHandler: NewExceptionHandler(rbc, rbp, kfExc, udp),
	}
}
