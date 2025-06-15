package twamp_test

import (
	"testing"
	"time"

	"github.com/snormore/twamp"
	"github.com/stretchr/testify/require"
)

func TestProbeSummary_Update(t *testing.T) {
	s := &twamp.ProbeSummary{}

	s.Update(12*time.Millisecond, true)
	s.Update(15*time.Millisecond, true)
	s.Update(11*time.Millisecond, true)
	s.Update(0, false) // Lost packet

	require.Equal(t, 3, s.Count)
	require.Equal(t, 1, s.Lost)
	require.Equal(t, 11*time.Millisecond, s.MinRTT)
	require.Equal(t, 15*time.Millisecond, s.MaxRTT)
	require.Equal(t, time.Duration(12+15+11)*time.Millisecond, s.TotalRTT)
	require.True(t, s.Jitter > 0)
	require.InDelta(t, 13*time.Millisecond, float64(s.AvgRTT()), float64(time.Millisecond))
}

func TestProbeSummary_InitialRTT(t *testing.T) {
	s := &twamp.ProbeSummary{}
	s.Update(9*time.Millisecond, true)

	require.Equal(t, 9*time.Millisecond, s.MinRTT)
	require.Equal(t, 9*time.Millisecond, s.MaxRTT)
	require.Equal(t, 9*time.Millisecond, s.AvgRTT())
	require.Equal(t, time.Duration(0), s.Jitter)
	require.Equal(t, 1, s.Count)
	require.Equal(t, 0, s.Lost)
}

func TestProbeSummary_ZeroRTT(t *testing.T) {
	s := &twamp.ProbeSummary{}
	s.Update(0, true)

	require.Equal(t, 0*time.Millisecond, s.MinRTT)
	require.Equal(t, 0*time.Millisecond, s.MaxRTT)
	require.Equal(t, 0*time.Millisecond, s.AvgRTT())
	require.Equal(t, 1, s.Count)
	require.Equal(t, 0, s.Lost)
}

func TestProbeSummary_MixedSequence(t *testing.T) {
	s := &twamp.ProbeSummary{}
	s.Update(10*time.Millisecond, true)
	s.Update(0, false)
	s.Update(15*time.Millisecond, true)
	s.Update(0, false)
	s.Update(12*time.Millisecond, true)

	require.Equal(t, 3, s.Count)
	require.Equal(t, 2, s.Lost)
	require.Equal(t, 10*time.Millisecond, s.MinRTT)
	require.Equal(t, 15*time.Millisecond, s.MaxRTT)
	require.Equal(t, time.Duration(10+15+12)*time.Millisecond, s.TotalRTT)
	require.True(t, s.Jitter > 0)
	require.InDelta(t, 12*time.Millisecond, float64(s.AvgRTT()), float64(time.Millisecond))
}
