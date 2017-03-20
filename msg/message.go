package msg

import (
	"goio/queue"
)

type Message interface {
	IOMessage
	Type() uint8
	SetChannel(Channel)
	Channel() Channel
}

type IOMessage interface {
	Encode(b *queue.IOBuffer) error
	Decode(b *queue.IOBuffer, len int32) error
}
