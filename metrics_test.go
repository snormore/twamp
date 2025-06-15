package twamp

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestMetricsCollector_Probes(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollector()
	mc.Register(reg)

	mc.ProbesSent.Inc()
	mc.ProbesReceived.Add(2)
	mc.PacketsDropped.Inc()
	mc.RTT.Observe(0.1)

	require.Equal(t, float64(1), testutil.ToFloat64(mc.ProbesSent))
	require.Equal(t, float64(2), testutil.ToFloat64(mc.ProbesReceived))
	require.Equal(t, float64(1), testutil.ToFloat64(mc.PacketsDropped))
	count := testutil.CollectAndCount(mc.RTT)
	require.GreaterOrEqual(t, count, 1)
}
