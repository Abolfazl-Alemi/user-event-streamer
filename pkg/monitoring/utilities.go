package monitoring

func (pr *ProMetrics) GaugeMetricIncr(labels map[string]string) {
	pr.EventGauge.With(labels).Inc()
}
