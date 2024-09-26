package streamer

import (
	"user-event-streamer/pkg/kafka"
	"user-event-streamer/pkg/monitoring"
	"user-event-streamer/pkg/rabbit"
	"user-event-streamer/pkg/udp"
)

type Streamer struct {
	udp *udp.Server
	kf  *kafka.Writer
	pr  *monitoring.ProMetrics
	rbc *rabbit.Consumer
	rbp *rabbit.Publisher
}

func NewStreamer(udp *udp.Server, kf *kafka.Writer, pr *monitoring.ProMetrics) *Streamer {
	return &Streamer{
		udp: udp,
		kf:  kf,
		pr:  pr,
	}
}
