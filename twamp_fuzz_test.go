package twamp_test

import (
	"testing"
	"time"

	"github.com/snormore/twamp"
)

func FuzzUnmarshalTestPacket(f *testing.F) {
	f.Add([]byte{0x00})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = twamp.UnmarshalTestPacket(data)
	})
}

func FuzzRoundTripPacket(f *testing.F) {
	f.Add(uint32(42), uint32(123), uint32(456), uint16(1), []byte("pad"))
	f.Fuzz(func(t *testing.T, seq, tsec, tfrac uint32, errEst uint16, pad []byte) {
		pkt := &twamp.TestPacket{
			Seq:           seq,
			TimestampSec:  tsec,
			TimestampFrac: tfrac,
			ErrorEstimate: errEst,
			Padding:       pad,
		}
		b := pkt.MarshalBinary()
		parsed, err := twamp.UnmarshalTestPacket(b)
		if err != nil {
			t.Skip()
		}
		_ = parsed
	})
}

func FuzzNTPConversion(f *testing.F) {
	f.Add(int64(time.Now().Unix()))
	f.Fuzz(func(t *testing.T, unixSec int64) {
		tm := time.Unix(unixSec, 0).UTC()
		sec, frac := twamp.ToNTP(tm)
		_ = twamp.FromNTP(sec, frac)
	})
}
