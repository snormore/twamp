package twamp

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

type ProbeResult struct {
	Seq      uint32        `json:"seq"`
	RTT      time.Duration `json:"rtt"`
	Sent     time.Time     `json:"sent"`
	Received time.Time     `json:"received"`
}

type Sender interface {
	SendProbe() (*ProbeResult, error)
	SendProbeWithPadding(paddingLen int) (*ProbeResult, error)
	Close() error
	Summary() *ProbeSummary
	ResetSummary()
}

type lightSender struct {
	log     *slog.Logger
	conn    *net.UDPConn
	remote  *net.UDPAddr
	timeout time.Duration
	seq     uint32
	lock    sync.Mutex
	metrics *MetricsCollector
	summary *ProbeSummary
}

func NewSender(local, remote string, timeout time.Duration, log *slog.Logger, metrics *MetricsCollector) (Sender, error) {
	if log == nil {
		log = slog.Default()
	}
	log = log.With("local", local, "remote", remote)

	localAddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		return nil, fmt.Errorf("resolving udp address: %w", err)
	}
	raddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return nil, fmt.Errorf("resolving udp address: %w", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("listening on udp address: %w", err)
	}

	return &lightSender{
		log:     log,
		conn:    conn,
		remote:  raddr,
		timeout: timeout,
		metrics: metrics,
		summary: &ProbeSummary{},
	}, nil
}

func (s *lightSender) Close() error {
	return s.conn.Close()
}

func (s *lightSender) Summary() *ProbeSummary {
	return s.summary
}

func (s *lightSender) ResetSummary() {
	s.summary = &ProbeSummary{}
}

func (s *lightSender) SendProbe() (*ProbeResult, error) {
	res, err := s.SendProbeWithPadding(0)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *lightSender) SendProbeWithPadding(paddingLen int) (*ProbeResult, error) {
	if s.metrics != nil {
		s.metrics.ProbesSent.Inc()
	}

	start := time.Now()

	s.lock.Lock()
	seq := s.seq
	s.seq++
	s.lock.Unlock()

	s.log.Debug("Sending probe", "seq", seq, "padding", paddingLen)

	sentTime := time.Now().UTC()
	sec, frac := ToNTP(sentTime)
	pkt := &TestPacket{
		Seq:           seq,
		TimestampSec:  sec,
		TimestampFrac: frac,
		ErrorEstimate: 1,
		MBZ:           0,
		Padding:       make([]byte, paddingLen),
	}
	data := pkt.MarshalBinary()
	_, err := s.conn.WriteToUDP(data, s.remote)
	if err != nil {
		s.summary.Update(0, false)
		return nil, fmt.Errorf("writing to udp: %w", err)
	}

	if err := s.conn.SetReadDeadline(time.Now().Add(s.timeout)); err != nil {
		s.summary.Update(0, false)
		return nil, fmt.Errorf("setting read deadline: %w", err)
	}
	buf := make([]byte, 2048)
	n, _, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		s.summary.Update(0, false)
		return nil, fmt.Errorf("reading from udp: %w", err)
	}
	recvTime := time.Now().UTC()
	reply, err := UnmarshalTestPacket(buf[:n])
	if err != nil {
		s.summary.Update(0, false)
		return nil, fmt.Errorf("unmarshalling reply: %w", err)
	}
	result := &ProbeResult{
		Seq:      reply.Seq,
		RTT:      recvTime.Sub(sentTime),
		Sent:     sentTime,
		Received: recvTime,
	}

	s.summary.Update(result.RTT, true)

	s.log.Debug("Received reply", "seq", reply.Seq, "rtt", result.RTT)

	if s.metrics != nil {
		s.metrics.RTT.Observe(time.Since(start).Seconds())
	}
	return result, nil
}
