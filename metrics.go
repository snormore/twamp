package twamp

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsCollector holds Prometheus metrics for TWAMP traffic.
type MetricsCollector struct {
	ProbesSent     prometheus.Counter
	ProbesReceived prometheus.Counter
	RTT            prometheus.Histogram
	PacketsDropped prometheus.Counter
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		ProbesSent: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "twamp",
			Name:      "probes_sent_total",
			Help:      "Total number of TWAMP probes sent.",
		}),
		ProbesReceived: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "twamp",
			Name:      "probes_received_total",
			Help:      "Total number of TWAMP probes received.",
		}),
		RTT: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "twamp",
			Name:      "rtt_seconds",
			Help:      "Round-trip time measured by TWAMP probes.",
			Buckets:   prometheus.ExponentialBuckets(1e-6, 2, 15),
		}),
		PacketsDropped: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "twamp",
			Name:      "packets_dropped_total",
			Help:      "Total number of malformed or dropped TWAMP packets.",
		}),
	}
}

func (m *MetricsCollector) Register(reg prometheus.Registerer) {
	reg.MustRegister(m.ProbesSent, m.ProbesReceived, m.RTT, m.PacketsDropped)
}
