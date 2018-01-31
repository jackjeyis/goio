package protocol

import (
	"errors"
	"goio/msg"
	"goio/queue"
)

type PingHeader struct {
	Op  byte
	Len int32
}

func (p *PingHeader) Decode(buf *queue.IOBuffer) error {
	temp := buf.Read(uint64(1))
	p.Op = temp[0] & 0xf0
	p.Len = decodeLength(buf)
	return nil
}

func (p *PingHeader) Encode(buf *queue.ByteBuffer) {
	p.Op = p.Op << 4
	buf.WriteByte(p.Op)
	encodeLength(buf, p.Len)
}

type PingPackage struct {
	id int
	ch msg.Channel
	PingHeader
	Body []byte
}

func (p *PingPackage) Type() uint8 {
	return uint8(p.Op)
}

func (p *PingPackage) SetChannel(ch msg.Channel) {
	p.ch = ch
}

func (p *PingPackage) Channel() msg.Channel {
	return p.ch
}

func (p *PingPackage) SetHandlerId(id int) {
	p.id = id
}

func (p *PingPackage) HandlerId() int {
	return p.id
}

func (p *PingPackage) Decode(buf *queue.IOBuffer, rlen int32) error {
	p.Body = buf.Read(uint64(p.Len))
	return nil
}

func (p *PingPackage) Encode(buf *queue.IOBuffer) error {
	b := queue.Get()
	p.PingHeader.Len = int32(len(p.Body))
	p.PingHeader.Encode(b)
	buf.Write(b.Bytes())
	buf.Write(p.Body)
	queue.Put(b)
	return nil
}

type PingProtocol struct {
	PingHeader
	rlen int32
}

func (p *PingProtocol) Encode(msg msg.Message, buf *queue.IOBuffer) error {
	return msg.Encode(buf)
}

func (p *PingProtocol) Decode(buf *queue.IOBuffer) (msg.Message, error) {
	var (
		cnt  = 2
		rlen int32
		err  error
	)
	if buf.Init() {
		for {
			if cnt > 5 {
				return nil, errors.New("error header size")
			}

			if buf.GetReadSize() < uint64(cnt+1) {
				return nil, HeaderErr
			}

			if buf.Byte(uint64(cnt)) >= 0x80 {
				cnt += 1
			} else {
				break
			}
		}
	}

	if buf.Init() {
		err = p.PingHeader.Decode(buf)
		buf.SetInit(false)
	} else {
		rlen = p.rlen
	}

	if uint64(rlen) > buf.GetReadSize() {
		p.rlen = p.PingHeader.Len
		return nil, BodyErr
	}

	if err != nil {
		return nil, err
	}

	msg := &PingPackage{}
	msg.PingHeader = p.PingHeader
	return msg, msg.Decode(buf, p.PingHeader.Len)
}
