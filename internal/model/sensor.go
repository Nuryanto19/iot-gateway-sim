package model

import (
	"encoding/binary"
	"math"
)

type SensorData struct {
	DeviceID uint32
	Value    float32
	Protocol string
}

func (s *SensorData) Pack() []byte {
	buf := make([]byte, 8)

	binary.BigEndian.PutUint32(buf[:4], s.DeviceID)

	bits := math.Float32bits(s.Value)
	binary.BigEndian.PutUint32(buf[4:], bits)

	return buf
}

func Unpack(b []byte) SensorData {
	return SensorData{
		DeviceID: binary.BigEndian.Uint32(b[:4]),
		Value:    math.Float32frombits(binary.BigEndian.Uint32(b[4:])),
	}
}
