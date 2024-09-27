package streamer

import (
	"encoding/json"
	"go.uber.org/zap"
	"sync"
	"user-event-streamer/internal/constants"
	"user-event-streamer/internal/models"
	"user-event-streamer/pkg/logger"
)

func (st *Streamer) RunPipeline(wg *sync.WaitGroup) {
	for event := range st.udp.UdpListenerChan {
		st.pr.GaugeMetricIncr(map[string]string{constants.GaugeLabelType: constants.GaugeValueTypeEvent, constants.GaugeLabelStatus: constants.GaugeValueStatusTotal})
		wg.Add(1)
		go func(data []byte) {
			defer wg.Done()
			var eData models.GA4Event
			err := json.Unmarshal(data, &eData)
			if err != nil {
				logger.Zap.Error("cannot unmarshal the GA4 event data", zap.Error(err), zap.String("event", string(data)))
				st.pr.GaugeMetricIncr(map[string]string{constants.GaugeLabelType: constants.GaugeValueTypeEvent, constants.GaugeLabelStatus: constants.GaugeValueStatusInvalid})
				if err := st.exceptionHandler.sendInvalidDataToKafka(data); err != nil {
					logger.Zap.Error("cannot publish the invalid event into kafka", zap.Error(err), zap.String("event", string(data)), zap.String("operation", "exception_handler"))
				}
				return
			}

			err = st.kf.Publish(data)
			if err != nil {
				logger.Zap.Error("cannot publish the event into kafka", zap.Error(err), zap.String("event", string(data)))
				if err := st.exceptionHandler.sendUnpublishedDataToRabbit(data); err != nil {
					logger.Zap.Error("cannot publish the exception event into rabbitMQ", zap.Error(err), zap.String("event", string(data)), zap.String("operation", "exception_handler"))
				}
			}
			st.pr.GaugeMetricIncr(map[string]string{constants.GaugeLabelType: constants.GaugeValueTypeEvent, constants.GaugeLabelStatus: constants.GaugeValueStatusValid})
		}(event)
	}
}
