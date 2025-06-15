package twamp_test

import (
	"testing"
	"time"

	"github.com/snormore/twamp"
	"github.com/stretchr/testify/require"
)

func TestNTPConversion_Precision(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	sec, frac := twamp.ToNTP(now)
	back := twamp.FromNTP(sec, frac)
	delta := now.Sub(back)
	if delta < 0 {
		delta = -delta
	}
	require.Less(t, delta, 500*time.Microsecond)
}
