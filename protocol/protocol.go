package protocol

import (
	"goio/msg"
	"goio/queue"
)

type Protocol interface {
	Encode(msg msg.Message, buf *queue.IOBuffer) error
	Decode(buf *queue.IOBuffer) (msg.Message, error)
}
