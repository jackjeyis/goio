package binary

import "goio/queue"

var BigEndian bigEndian

type bigEndian struct{}

func (bigEndian) Int16(b []byte) int16 {
	return int16(b[1]) | int16(b[0])<<8
}

func (bigEndian) PutInt16(b *queue.ByteBuffer, v int16) {
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}

func (bigEndian) Int32(b []byte) int32 {
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

func (bigEndian) PutInt32(b *queue.ByteBuffer, v int32) {
	b.WriteByte(byte(v >> 24))
	b.WriteByte(byte(v >> 16))
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}
