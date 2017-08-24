package protocol

import (
	"goio/msg"
	"goio/queue"
)

type WebSocketHeader struct {
	fin     bool   // 1 over 0 continue
	rsv     byte   // default 0 observe
	opcode  byte   // 0 continution frame 1 text frame 2 binary frame 8 close frame 9 ping frame 10 pong frame
	mask    bool   // client to server must 1 else close | server to client must 0 else close
	len     byte   // payload len 126 |2 bytes length| 127 |8 bytes length|
	maskkey uint32 // mask number when mask is 1
}

type WebSocketPackage struct {
	id int
	ch msg.Channel
	WebSocketHeader
	Body []byte
}

func (w *WebSocketPackage) Type() uint8 {
	return uint8(0)
}

func (w *WebSocketPackage) SetChannel(ch msg.Channel) {
	w.ch = ch
}

func (w *WebSocketPackage) Channel() msg.Channel {
	return w.ch
}

func (w *WebSocketPackage) SetHandlerId(id int) {
	w.id = id
}

func (w *WebSocketPackage) HandlerId() int {
	return w.id
}

func (w *WebSocketPackage) Decode(b *queue.IOBuffer, rlen int32) error {
	return nil
}

func (w *WebSocketPackage) Encode(b *queue.IOBuffer) error {
	return nil
}

type WebSocketProtocol struct {
}

func (w *WebSocketProtocol) Encode(msg msg.Message, buf *queue.IOBuffer) error {
	return nil
}

func (w *WebSocketProtocol) Decode(buf *queue.IOBuffer) (msg.Message, error) {
	return nil, nil
}
