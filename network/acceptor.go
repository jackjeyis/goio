package network

import (
	"goio/logger"
	"goio/msg"
	"goio/protocol"
	"goio/service"
	"goio/netpoll"
	"goio/pool"
	"net"
	"os"
	"os/exec"
	"runtime"
	"github.com/valyala/fasthttp/reuseport"
	"time"
	"fmt"
)

type Acceptor struct {
	io_service  *service.IOService
	service     msg.Service
	addr        string
	ln          *net.TCPListener
	proto       protocol.Protocol
	quit        chan struct{}
	freeChannel *ServiceChannel
}

func NewAcceptor(io_srv *service.IOService, srv msg.Service, address string, proto protocol.Protocol) *Acceptor {
	return &Acceptor{
		io_service: io_srv,
		service:    srv,
		addr:       address,
		proto:      proto,
		quit:       make(chan struct{}),
	}
}
var child = true
func (a *Acceptor) GetListener() (*net.TCPListener,error) {
	if child {
	children := make([]*exec.Cmd,runtime.NumCPU()/100)
	for i := range children {
		children[i] = exec.Command("taskset","-c",fmt.Sprintf("%d",i),os.Args[0],"-c","conf.toml","-child")
		children[i].Stdout = os.Stdout
		children[i].Stderr = os.Stderr
		if err := children[i].Start(); err != nil {
			logger.Error("start err %v",err)
		}
	}
	
	for _, ch := range children {
		if err := ch.Wait(); err != nil {
			logger.Error("wait err %v",err)
		}
	}
	os.Exit(0)
	panic("unreachable")
	}
	ln, err := reuseport.Listen("tcp4",a.addr)
	if err != nil {
		logger.Error("listen err %v",err)
	}
	return ln.(*net.TCPListener),nil
}

func (a *Acceptor) Start() {
	defer a.ln.Close()
	netaddr, err := net.ResolveTCPAddr("tcp", a.addr)
	if err != nil {
		logger.Error("net.ResolveTCPAddr fail %v", netaddr)
		return
	}

	logger.Info("netaddr %v", netaddr)
	a.ln, err = /*net.ListenTCP("tcp", netaddr)*/ a.GetListener()
	if err != nil {
		logger.Error("net.ListenTCP fail %v",err)
		return
	}
	poller, _ := netpoll.New(nil)

	wp := &pool.WorkerPool{
		MaxWorkersCount: 1000,
	}

	handle := func(conn net.Conn){
		desc := netpoll.Must(netpoll.HandleRead(conn))
		poller.Start(desc,func(e netpoll.Event){
			if e & (netpoll.EventReadHup|netpoll.EventHup) != 0 {
				poller.Stop(desc)
				return
			}
			wp.ScheduleTimeout(time.Millisecond,func(){
				buf := make([]byte,128)
				if n, err := conn.Read(buf); err != nil {
					poller.Stop(desc)
					logger.Error("n %v,err %v",n,err)
				}	
			})	
		})
	}

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		a.ln, netpoll.EventRead/*|netpoll.EventOneShot*/,
	))

	//accept := make(chan error, 1)
	poller.Start(acceptDesc, func(e netpoll.Event) {
		wp.ScheduleTimeout(time.Millisecond, func() {
			conn, err := a.ln.AcceptTCP()
			logger.Info("connected %v, err %v", conn, err)
			handle(conn)
		})
		poller.Resume(acceptDesc)
	})
	
	/*for i := 0; i < runtime.NumCPU(); i++ {
		go acceptTCP(a)
	}*/
}

func acceptTCP(a *Acceptor) {

	var (
		conn  *net.TCPConn
		delay time.Duration
		err   error
	)

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
		//conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		sch := a.allocChannel()
		InitChannel(sch, a, conn)
		sch.Start()
		logger.Info("accept new conn %v", conn)
	}
}

func (a *Acceptor) Stop() {
	a.ln.Close()
	close(a.quit)
	Close()
}

func (a *Acceptor) allocChannel() *ServiceChannel {
	ch := a.freeChannel
	if ch == nil {
		ch = new(ServiceChannel)
	} else {
		a.freeChannel = ch.next
		*ch = ServiceChannel{}
	}
	return ch
	//return &ServiceChannel{}
}
