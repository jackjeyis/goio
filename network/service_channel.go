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
	queue      *queue.Sequence
}

func NewServiceChannel(conn *net.TCPConn, proto protocol.Protocol,
	io_srv *service.IOService, srv msg.Service,
	q chan struct{}) *ServiceChannel {
	return &ServiceChannel{
		in:         queue.NewIOBuffer(true),
		out:        queue.NewIOBuffer(false),
		conn:       conn,
		close:      make(chan struct{}),
		attrs:      make(map[string]string),
		proto:      proto,
		io_service: io_srv,
		service:    srv,
		write:      make(chan bool, 10000),
		quit:       q,
		queue:      queue.NewSequence(1024 * 1024),
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
	go s.OnWrite()
}

func (s *ServiceChannel) OnRead() {
	var (
		err error
		wt  uint64
		n   int
	)
	defer func() {
		s.conn.CloseRead()
		s.in.Reset()
		s.OnClose()
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
			break L
		}
		s.in.Produce(uint64(n))
		//logger.Info("n bytes %v,wt %v,rt %v,size %v,read size %v", n, s.in.GetWrite(), s.in.GetRead(), s.in.Len(), s.in.GetReadSize())
		err = s.DecodeMessage()
		if err != nil {
			logger.Info("Decode error %v", err)
			if err == protocol.HeaderErr || err == protocol.BodyErr {
				continue
			}

			break L
		}
	}
}

func (s *ServiceChannel) OnWrite() {
	var (
		n int
	)
	defer func() {
		s.out.Reset()
		close(s.write)
	}()

	for {
		select {
		case <-s.quit:
			s.conn.Close()
			return
		case <-s.close:
			if s.out.GetReadSize() > 0 {
				continue
			} else {
				s.conn.CloseWrite()
				return
			}

		default:
			/*if s.out.GetReadSize() > 0 {
				//logger.Info("wt %v,rt %v,cap %v,size %v", s.out.GetWrite(), s.out.GetRead(), s.out.Len(), s.out.GetReadSize())
				n, err = s.conn.Write(s.out.Buffer()[s.out.GetRead():s.out.GetWrite()])
				if n < 0 || err != nil {
					s.conn.CloseWrite()
					return
				}
				s.out.Consume(uint64(n))
			}
			*/
			m, err := s.queue.Get()
			if err != nil {
				continue
			}
			if msg, ok := m.(msg.Message); ok {
				if err := s.proto.Encode(msg, s.out); err != nil {
					logger.Error("s.protocol.Encode error %v", err)
					s.conn.Close()
					return
				}
				n, err = s.conn.Write(s.out.Buffer()[s.out.GetRead():s.out.GetWrite()])
				if n < 0 || err != nil {
					s.conn.CloseWrite()
					return
				}
				s.out.Consume(uint64(n))
				logger.Info("n %v,wt %v,rt %v,cap %v,size %v", n, s.out.GetWrite(), s.out.GetRead(), s.out.Len(), s.out.GetReadSize())
			}
		}

	}
}

func (s *ServiceChannel) DecodeMessage() error {
	var (
		m msg.Message
		e error
	)
	for s.in.GetReadSize() > 0 && e == nil {
		m, e = s.proto.Decode(s.in)
		if e != nil {
			return e
		}
		m.SetChannel(s)
		s.io_service.GetServiceStage().Send(m)
		s.in.SetInit(true)
	}
	return nil
}

func (s *ServiceChannel) EncodeMessage(msg msg.Message) {
	//logger.Info("msg %v", msg)
	/*if err := s.proto.Encode(msg, s.out); err != nil {
		logger.Error("s.protocol.Encode error %v", err)
		s.conn.Close()
		return
	}*/
	s.queue.Put(msg)
}

func (s *ServiceChannel) Serve(msg msg.Message) {
	s.service.Serve(msg)
}

func (s *ServiceChannel) OnClose() {
	close(s.close)
}
