package network

import (
	"goio/logger"
	"goio/msg"
	"goio/protocol"
	"goio/service"
	"net"
	"sync"
	"time"
)

type Acceptor struct {
	io_service *service.IOService
	service    msg.Service
	addr       string
	ln         *net.TCPListener
	proto      protocol.Protocol
	quit       chan struct{}
	wg         *sync.WaitGroup
}

func NewAcceptor(io_srv *service.IOService, srv msg.Service, address string, proto protocol.Protocol) *Acceptor {
	return &Acceptor{
		io_service: io_srv,
		service:    srv,
		addr:       address,
		proto:      proto,
		quit:       make(chan struct{}),
		wg:         &sync.WaitGroup{},
	}
}

func (a *Acceptor) Start() {
	defer a.ln.Close()
	netaddr, err := net.ResolveTCPAddr("tcp", a.addr)
	if err != nil {
		logger.Error("net.ResolveTCPAddr fail %v", netaddr)
		return
	}

	logger.Info("netaddr %v", netaddr)
	a.ln, err = net.ListenTCP("tcp", netaddr)
	if err != nil {
		logger.Error("net.ListenTCP fail")
		return
	}

	//a.ln.SetDeadline(time.Now().Add(10 * time.Second))

	var (
		conn  *net.TCPConn
		delay time.Duration
	)

	//var JobQueue chan msg.ProtocolMessage

	for {
		conn, err = a.ln.AcceptTCP()
		if err != nil {
			nerr, ok := err.(net.Error)
			if ok {
				if nerr.Timeout() {
					logger.Error("AcceptTCP timeout")
					return
				}
				if nerr.Temporary() {
					if delay == 0 {
						delay = 5 * time.Millisecond
					} else {
						delay *= 2
					}
					if max := 1 * time.Second; delay > max {
						delay = max
					}

					logger.Error("AcceptTCP Temporary error: %v,retry after : %v", err, delay)
					time.Sleep(delay)
					continue
				}
			}
			return
		}

		//conn.SetReadDeadline(time.Now().Add(6 * time.Second))
		serviceChannel := NewServiceChannel(
			conn,
			a.proto,
			a.io_service,
			a.service,
			a.quit,
		)

		a.wg.Add(1)
		go func(serviceChannel *ServiceChannel) {
			serviceChannel.Start()
			a.wg.Done()
		}(serviceChannel)

		/*	msg := &ProtocolMessage{
				channel: serviceChannel,
			}

			JobQueue <- msg

			dispacther := &Dispathcher{
				queue: JobQueue,
			}

			dispacther.Run()
		*/
	}
}

func (a *Acceptor) Stop() {
	a.ln.Close()
	close(a.quit)
	a.wg.Wait()
}
