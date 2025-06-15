package twamp

import (
	"encoding/binary"
	"errors"
	"slices"
)

type TestPacket struct {
	Seq               uint32
	TimestampSec      uint32
	TimestampFrac     uint32
	ErrorEstimate     uint16
	MBZ               uint16
	RecvTimestampSec  uint32
	RecvTimestampFrac uint32
	Padding           []byte
}

func (p *TestPacket) MarshalBinary() []byte {
	buf := make([]byte, 24+len(p.Padding))
	binary.BigEndian.PutUint32(buf[0:4], p.Seq)
	binary.BigEndian.PutUint32(buf[4:8], p.TimestampSec)
	binary.BigEndian.PutUint32(buf[8:12], p.TimestampFrac)
	binary.BigEndian.PutUint16(buf[12:14], p.ErrorEstimate)
	binary.BigEndian.PutUint16(buf[14:16], p.MBZ)
	binary.BigEndian.PutUint32(buf[16:20], p.RecvTimestampSec)
	binary.BigEndian.PutUint32(buf[20:24], p.RecvTimestampFrac)
	copy(buf[24:], p.Padding)
	return buf
}

func UnmarshalTestPacket(buf []byte) (*TestPacket, error) {
	if len(buf) < 24 {
		return nil, errors.New("TWAMP-Test packet too short")
	}
	return &TestPacket{
		Seq:               binary.BigEndian.Uint32(buf[0:4]),
		TimestampSec:      binary.BigEndian.Uint32(buf[4:8]),
		TimestampFrac:     binary.BigEndian.Uint32(buf[8:12]),
		ErrorEstimate:     binary.BigEndian.Uint16(buf[12:14]),
		MBZ:               binary.BigEndian.Uint16(buf[14:16]),
		RecvTimestampSec:  binary.BigEndian.Uint32(buf[16:20]),
		RecvTimestampFrac: binary.BigEndian.Uint32(buf[20:24]),
		Padding:           slices.Clone(buf[24:]),
	}, nil
}
