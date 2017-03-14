package network

import (
	"io"
	"net"

	"goio/logger"
	"goio/msg"
	"goio/protocol"
	"goio/queue"
	"goio/service"
)

const (
	per_read_size    = 4096
	decode_watermask = 4096
)

var (
	m   msg.Message
	err error
)

type ServiceChannel struct {
	close      chan struct{}
	quit       chan struct{}
	write      chan bool
	in         *queue.IOBuffer
	out        *queue.IOBuffer
	conn       *net.TCPConn
	proto      protocol.Protocol
	io_service *service.IOService
	service    msg.Service
	attrs      map[string]string
}

func NewServiceChannel(conn *net.TCPConn, proto protocol.Protocol, io_srv *service.IOService, srv msg.Service, q chan struct{}) *ServiceChannel {
	return &ServiceChannel{
		in:         queue.NewIOBuffer(),
		out:        queue.NewIOBuffer(),
		conn:       conn,
		close:      make(chan struct{}),
		attrs:      make(map[string]string),
		proto:      proto,
		io_service: io_srv,
		service:    srv,
		write:      make(chan bool),
		quit:       q,
	}
}

func (s *ServiceChannel) SetAttr(key, value string) {
	s.attrs[key] = value
}

func (s *ServiceChannel) GetAttr(key string) string {
	if attr, ok := s.attrs[key]; ok {
		return attr
	}
	return ""
}

func (s *ServiceChannel) Start() {
	go s.OnRead()
	s.OnWrite()
}

func (s *ServiceChannel) OnRead() {
	var (
		wt uint64
		n  int
	)
	defer func() {
		s.conn.CloseRead()
		s.in.Reset()
	}()
L:
	for {
		wt, err = s.in.EnsureWrite(per_read_size)
		if err != nil {
			logger.Error("EnsureWrite %v", err)
			return
		}

		n, err = s.conn.Read(s.in.Buffer()[wt:])

		if err != nil {
			if err == io.EOF {
				logger.Warn("ServiceChannel.OnRead Closed by peer %v", err)
			} else {
				logger.Error("ServiceChannel.OnRead Error %v", err)
			}
			s.OnClose()
			break L
		}
		s.in.Produce(uint64(n))
		err = s.DecodeMessage()
		if err != nil {
			return
		}
	}
	s.DecodeMessage()
}

func (s *ServiceChannel) OnWrite() {
	var n int
	defer func() {
		s.out.Reset()
	}()

	for {
		select {
		case <-s.write:
			n, err = s.conn.Write(s.out.Buffer()[s.out.GetRead():s.out.GetWrite()])
			if n <= 0 || err != nil {
				s.conn.CloseWrite()
				return
			}
			s.out.Consume(uint64(n))
		case <-s.quit:
			if s.out.GetReadSize() > 0 {
				continue
			} else {
				s.conn.Close()
				return
			}
		case <-s.close:
			if s.out.GetReadSize() > 0 {
				continue
			} else {
				s.conn.CloseWrite()
				return
			}
		}
	}
}

func (s *ServiceChannel) DecodeMessage() (e error) {
	for s.in.GetReadSize() > 0 && e == nil {
		m, e = s.proto.Decode(s.in)
		if e != nil {
			logger.Info("Decode error %v", e)
			s.OnClose()
			return e
		}
		m.Protocol(msg.NewProtocolMessage(s))
		s.io_service.GetServiceStage().Send(m)
	}
	return nil
}

func (s *ServiceChannel) EncodeMessage(msg msg.Message) {
	if err = s.proto.Encode(msg, s.out); err != nil {
		logger.Error("s.protocol.Encode error %v", err)
	}
	s.write <- true
}

func (s *ServiceChannel) Service() msg.Service {
	return s.service
}

func (s *ServiceChannel) OnClose() {
	close(s.close)
}
