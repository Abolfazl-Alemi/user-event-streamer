package monitoring

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strings"
)

type ProMetrics struct {
	TimeToProcess *prometheus.SummaryVec
	EventGauge    *prometheus.GaugeVec
}

func NewPrometheusMonitoring(reg prometheus.Registerer, summaryLabels []string, gaugeLabels []string) *ProMetrics {
	// Open the go.mod file
	file, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing go.mod:", err)
		}
	}(file)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var moduleName string = "myProject"
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			// The module name is the second part of the line
			moduleName = strings.Fields(line)[1]
			fmt.Println("Module name:", moduleName)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading go.mod:", err)
	}
	metrics := &ProMetrics{
		TimeToProcess: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:   moduleName,
				Name:        "response_time_summary",
				Help:        "the time taken to do an operation",
				ConstLabels: nil,
				Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				MaxAge:      0,
				AgeBuckets:  0,
				BufCap:      0,
			},
			summaryLabels,
		),
		EventGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "events_gauge",
				Help:      "count of events group by status",
			},
			gaugeLabels,
		),
	}

	reg.MustRegister(metrics.TimeToProcess, metrics.EventGauge)

	return metrics
}
