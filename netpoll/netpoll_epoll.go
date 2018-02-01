package netpoll

import "goio/logger"

func New(c *Config) (Poller, error) {
	epoll, err := EpollCreate(nil)
	if err != nil {
		logger.Info("EpollCreate failed %v,config %v", err,c)
		return nil, err
	}
	return &poller{epoll},nil
}

type poller struct {
	*Epoll
}

func (p *poller) Start(desc *Desc, cb CallbackFn) error {
	return p.Add(desc.fd(), toEpollEvent(desc.event),
		func(ev EpollEvent) {
			var event Event
			if ev&EPOLLHUP != 0 {
				event |= EventHup
			}

			if ev&EPOLLRDHUP != 0 {
				event |= EventReadHup
			}

			if ev&EPOLLIN != 0 {
				event |= EventRead
			}

			if ev&EPOLLOUT != 0 {
				event |= EventWrite
			}

			if ev&EPOLLERR != 0 {
				event |= EventErr
			}

			if ev&_EPOLLCLOSED != 0 {
				event |= EventPollerClosed
			}
			cb(event)
		})
}

func (p *poller) Stop(desc *Desc) error {
	return p.Del(desc.fd())
}

func (p *poller) Resume(desc *Desc) error {
	return p.Mod(desc.fd(), toEpollEvent(desc.event))
}

func toEpollEvent(event Event) (ep EpollEvent) {
	if event&EventRead != 0 {
		ep |= EPOLLIN | EPOLLRDHUP
	}

	if event&EventWrite != 0 {
		ep |= EPOLLOUT
	}

	if event&EventOneShot != 0 {
		ep |= EPOLLONESHOT
	}

	if event&EventEdgeTriggered != 0 {
		ep |= EPOLLET
	}

	return ep
}
