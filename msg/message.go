package msg

import (
	"wolfhead/queue"
)

type Message interface {
	IOMessage
	Type() uint8
	Protocol(*ProtocolMessage)
	ProtoMsg() *ProtocolMessage
}

type IOMessage interface {
	Encode(b *queue.IOBuffer) error
	Decode(b *queue.IOBuffer, len int32) error
}
