package twamp_test

import (
	"testing"

	"github.com/snormore/twamp"
	"github.com/stretchr/testify/require"
)

func TestTestPacket_Roundtrip(t *testing.T) {
	original := &twamp.TestPacket{
		Seq:               42,
		TimestampSec:      0xDD000000,
		TimestampFrac:     0x12345678,
		ErrorEstimate:     0x0001,
		MBZ:               0,
		RecvTimestampSec:  0xAA000000,
		RecvTimestampFrac: 0x87654321,
		Padding:           []byte("abcdefg"),
	}
	dat := original.MarshalBinary()
	parsed, err := twamp.UnmarshalTestPacket(dat)
	require.NoError(t, err)
	checkRoundTripMatch(t, original, parsed)
}

func TestUnmarshal_InvalidLength(t *testing.T) {
	_, err := twamp.UnmarshalTestPacket([]byte{0x00})
	require.ErrorContains(t, err, "too short")
}

func checkRoundTripMatch(t *testing.T, original, parsed *twamp.TestPacket) {
	t.Helper()
	require.Equal(t, original.Seq, parsed.Seq)
	require.Equal(t, original.TimestampSec, parsed.TimestampSec)
	require.Equal(t, original.TimestampFrac, parsed.TimestampFrac)
	require.Equal(t, original.RecvTimestampSec, parsed.RecvTimestampSec)
	require.Equal(t, original.RecvTimestampFrac, parsed.RecvTimestampFrac)
	require.Equal(t, original.Padding, parsed.Padding)
}
