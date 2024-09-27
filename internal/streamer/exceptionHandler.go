package streamer

import (
	"sync"
	"user-event-streamer/pkg/kafka"
	"user-event-streamer/pkg/rabbit"
	"user-event-streamer/pkg/udp"
)

type ExceptionHandler struct {
	rbc   *rabbit.Consumer
	rbp   *rabbit.Publisher
	kfExc *kafka.Writer
	udp   *udp.Server
}

func NewExceptionHandler(rbc *rabbit.Consumer, rbp *rabbit.Publisher, kfExc *kafka.Writer, udp *udp.Server) *ExceptionHandler {
	return &ExceptionHandler{
		rbc:   rbc,
		rbp:   rbp,
		kfExc: kfExc,
		udp:   udp,
	}
}

func (ex *ExceptionHandler) sendUnpublishedDataToRabbit(data []byte) error {
	if err := ex.rbp.PublishMsg(data); err != nil {
		return err
	}
	return nil
}

func (ex *ExceptionHandler) sendInvalidDataToKafka(data []byte) error {
	if err := ex.kfExc.Publish(data); err != nil {
		return err
	}
	return nil
}

func (ex *ExceptionHandler) tryProcessUnpublishedData(wg *sync.WaitGroup) {
	for msg := range ex.rbc.ConsumerOutChan {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ex.udp.UdpListenerChan <- msg
		}()
	}
}
