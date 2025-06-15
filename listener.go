package twamp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"
)

type Listener interface {
	Run(ctx context.Context) error
	Close() error
	LocalAddr() net.Addr
}

type lightListener struct {
	log       *slog.Logger
	conn      *net.UDPConn
	reflector Reflector
	bufsize   int
}

func NewListener(addr string, bufsize int, log *slog.Logger, metrics *MetricsCollector) (Listener, error) {
	if log == nil {
		log = slog.Default()
	}

	reflector := NewReflector(log, metrics)

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("resolving udp address: %w", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("listening on udp address: %w", err)
	}

	log = log.With("addr", udpAddr)

	log.Info("Initialized TWAMP light listener")

	return &lightListener{
		log:       log,
		conn:      conn,
		reflector: reflector,
		bufsize:   bufsize,
	}, nil
}

func (s *lightListener) Run(ctx context.Context) error {
	s.log.Info("Running TWAMP light listener")

	buf := make([]byte, s.bufsize)
	for {
		select {
		case <-ctx.Done():
			s.log.Info("TWAMP light listener stopped by context")
			return nil
		default:
		}
		if err := s.conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
			return fmt.Errorf("setting read deadline: %w", err)
		}
		n, addr, err := s.conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if errors.Is(err, net.ErrClosed) {
				s.log.Info("UDP connection closed, stopping listener")
				return nil
			}
			s.log.Error("Error reading from UDP", "error", err)
			return fmt.Errorf("reading from udp: %w", err)
		}
		if n < 24 {
			s.log.Warn("Received malformed packet", "length", n)
			if r, ok := s.reflector.(*lightReflector); ok {
				handleMalformed(r.metrics)
			}
			continue
		}
		resp, err := s.reflector.HandleProbe(buf[:n], addr)
		if err != nil {
			s.log.Error("Error handling probe", "error", err)
			continue
		}
		if _, err := s.conn.WriteTo(resp, addr); err != nil {
			// Log warning but continue.
			s.log.Warn("writing response to UDP", "error", err)
		}
	}
}

func handleMalformed(metrics *MetricsCollector) {
	if metrics != nil {
		metrics.PacketsDropped.Inc()
	}
}

func (s *lightListener) Close() error {
	return s.conn.Close()
}

func (s *lightListener) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}
