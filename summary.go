package twamp

import (
	"sync"
	"time"
)

type ProbeSummary struct {
	mu       sync.Mutex
	Count    int
	Lost     int
	MinRTT   time.Duration
	MaxRTT   time.Duration
	TotalRTT time.Duration
	Jitter   time.Duration
	LastRTT  time.Duration
}

func (s *ProbeSummary) Update(rtt time.Duration, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !ok {
		s.Lost++
		return
	}
	s.Count++
	s.TotalRTT += rtt
	if s.MinRTT == 0 || rtt < s.MinRTT {
		s.MinRTT = rtt
	}
	if rtt > s.MaxRTT {
		s.MaxRTT = rtt
	}
	if s.LastRTT > 0 {
		delta := rtt - s.LastRTT
		if delta < 0 {
			delta = -delta
		}
		s.Jitter += delta
	}
	s.LastRTT = rtt
}

func (s *ProbeSummary) AvgRTT() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Count == 0 {
		return 0
	}
	return s.TotalRTT / time.Duration(s.Count)
}
