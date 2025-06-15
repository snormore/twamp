package twamp

import "time"

func ToNTP(t time.Time) (uint32, uint32) {
	secs := uint32(t.Unix() + 2208988800)
	nanos := uint64(t.Nanosecond())
	frac := uint32((nanos << 32) / 1e9)
	return secs, frac
}

func FromNTP(secs, frac uint32) time.Time {
	s := int64(secs) - 2208988800
	n := (uint64(frac) * 1e9) >> 32
	return time.Unix(s, int64(n)).UTC()
}
