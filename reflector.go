package twamp

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

type Reflector interface {
	HandleProbe([]byte, net.Addr) ([]byte, error)
}

// lightReflector implements Reflector
type lightReflector struct {
	log     *slog.Logger
	metrics *MetricsCollector
}

func NewReflector(log *slog.Logger, metrics *MetricsCollector) Reflector {
	if log == nil {
		log = slog.Default()
	}

	return &lightReflector{
		log:     log,
		metrics: metrics,
	}
}

func (r *lightReflector) HandleProbe(msg []byte, from net.Addr) ([]byte, error) {
	r.log.Debug("Received probe", "from", from)

	if r.metrics != nil {
		r.metrics.ProbesReceived.Inc()
	}

	pkt, err := UnmarshalTestPacket(msg)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling probe: %w", err)
	}
	sec, frac := ToNTP(time.Now().UTC())
	pkt.RecvTimestampSec = sec
	pkt.RecvTimestampFrac = frac

	r.log.Debug("Echoed response", "seq", pkt.Seq, "from", from)

	return pkt.MarshalBinary(), nil
}
