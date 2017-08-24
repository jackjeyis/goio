package protocol

import (
	"bytes"
	"errors"
	"goio/logger"
	"goio/msg"
	"goio/queue"
	"strconv"
)

type HttpHeader map[string]string

type HttpReq struct {
	id      int
	ch      msg.Channel
	Method  string
	Uri     string
	Version string
	Header  HttpHeader
	Body    []byte
}

func (h *HttpReq) Type() uint8 {
	return 0
}

func (h *HttpReq) SetChannel(ch msg.Channel) {
	h.ch = ch
}

func (h *HttpReq) Channel() msg.Channel {
	return h.ch
}

func (h *HttpReq) SetHandlerId(id int) {
	h.id = id
}

func (h *HttpReq) HandlerId() int {
	return h.id
}

func (h *HttpReq) Decode(b *queue.IOBuffer, rlen int32) error {
	var (
		line []byte
		pos  int
	)
	line, pos = b.ReadSlice('\n')
	if pos > 0 {
		s1 := bytes.IndexByte(line, ' ')
		s2 := bytes.IndexByte(line[s1+1:], ' ')
		s2 += s1 + 1
		if s1 < 0 || s2 < 0 {
			return errors.New("Bad Request!")
		}
		h.Method = string(line[:s1])
		h.Uri = string(line[s1+1 : s2])
		h.Version = string(line[s2+1:])

		buf := b.Buffer()[b.GetRead():b.GetWrite()]
		//logger.Info("line %s,method %s,uri %s ver %s,buffer %s", string(line), h.Method, h.Uri, h.Version, string(buf))
		h.Header = make(HttpHeader)
		for len(buf) > 0 {
			line, pos = b.ReadSlice('\n')
			if pos < 3 {
				break
			}
			colon := bytes.IndexByte(line, ':')
			h.Header[string(line[:colon])] = string(line[colon+1:])
			logger.Info("line %s,h.Header %v", string(line), h.Header)
			//return errors.New("some bug")
			buf = buf[pos+1:]
		}
	} else {
		return HeaderErr
	}
	clen := h.Header["Content-Length"]
	if len(clen) == 0 {
		return nil
	}
	n, err := strconv.ParseInt(clen, 10, 64)
	if err != nil {
		return err
	}
	h.Body = b.Read(uint64(n))
	return nil
}

func (h *HttpReq) Encode(b *queue.IOBuffer) error {
	return nil
}

type HttpProtocol struct {
	Req HttpReq
}

func (h *HttpProtocol) Encode(msg msg.Message, buf *queue.IOBuffer) error {
	return nil
}

func (h *HttpProtocol) Decode(buf *queue.IOBuffer) (msg.Message, error) {
	req := &HttpReq{}
	return req, req.Decode(buf, 0)
}
