package netpoll

import (
	"fmt"
	"sync"

	"golang.org/x/sys/unix"
)

type EpollEvent uint32

const (
	EPOLLIN      = unix.EPOLLIN
	EPOLLOUT     = unix.EPOLLOUT
	EPOLLRDHUP   = unix.EPOLLRDHUP
	EPOLLPRI     = unix.EPOLLPRI
	EPOLLERR     = unix.EPOLLERR
	EPOLLHUP     = unix.EPOLLHUP
	EPOLLET      = unix.EPOLLET
	EPOLLONESHOT = unix.EPOLLONESHOT
	_EPOLLCLOSED = 0x20
)

func (evt EpollEvent) String() (str string) {
	name := func(event EpollEvent, name string) {
		if evt&event == 0 {
			return
		}
		if str != "" {
			str += "|"
		}
		str += name
	}

	name(EPOLLIN, "EPOLLIN")
	name(EPOLLOUT, "EPOLLOUT")
	name(EPOLLRDHUP, "EPOLLRDHUP")
	name(EPOLLPRI, "EPOLLPRI")
	name(EPOLLERR, "EPOLLERR")
	name(EPOLLHUP, "EPOLLHUP")
	name(EPOLLET, "EPOLLET")
	name(EPOLLONESHOT, "EPOLLONESHOT")
	name(_EPOLLCLOSED, "_EPOLLCLOSED")

	return
}

type Epoll struct {
	lock sync.RWMutex

	fd      int
	eventfd int
	closed  bool
	done    chan struct{}

	callbacks map[int]func(EpollEvent)
}

func EpollCreate(waitErr func(error)) (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)

	if err != nil {
		return nil, err
	}

	r0, _, errno := unix.Syscall(unix.SYS_EVENTFD2, 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	eventFd := int(r0)

	err = unix.EpollCtl(fd, unix.EPOLL_CTL_ADD, eventFd, &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(eventFd),
	})

	if err != nil {
		unix.Close(fd)
		unix.Close(eventFd)
		return nil, err
	}

	ep := &Epoll{
		fd:        fd,
		eventfd:   eventFd,
		done:      make(chan struct{}),
		callbacks: make(map[int]func(EpollEvent)),
	}
	go ep.wait(waitErr)
	return ep, nil
}

var closeBytes = []byte{1, 0, 0, 0, 0, 0, 0, 0}

func (ep *Epoll) Close() (err error) {
	ep.lock.Lock()
	{
		if ep.closed {
			ep.lock.Unlock()
			return ErrClosed
		}
		ep.closed = true
		if _, err = unix.Write(ep.eventfd, closeBytes); err != nil {
			ep.lock.Unlock()
			return
		}
	}
	ep.lock.Unlock()
	<-ep.done
	if err = unix.Close(ep.eventfd); err != nil {
		return
	}
	ep.lock.Lock()
	callbacks := ep.callbacks
	ep.callbacks = nil
	ep.lock.Unlock()

	for _, cb := range callbacks {
		if cb != nil {
			cb(_EPOLLCLOSED)
		}
	}
	return
}

func (ep *Epoll) Add(fd int, events EpollEvent, cb func(EpollEvent)) (err error) {
	ev := &unix.EpollEvent{
		Events: uint32(events),
		Fd:     int32(fd),
	}
	ep.lock.Lock()
	defer ep.lock.Unlock()
	if ep.closed {
		return ErrClosed
	}

	if _, ok := ep.callbacks[fd]; ok {
		return ErrRegistered
	}
	ep.callbacks[fd] = cb
	return unix.EpollCtl(ep.fd, unix.EPOLL_CTL_ADD, fd, ev)
}

func (ep *Epoll) Del(fd int) (err error) {
	ep.lock.Lock()
	defer ep.lock.Unlock()
	if ep.closed {
		return ErrClosed
	}

	if _, ok := ep.callbacks[fd]; !ok {
		return ErrNotRegistered
	}

	delete(ep.callbacks, fd)
	return unix.EpollCtl(fd, unix.EPOLL_CTL_DEL, fd, nil)
}

func (ep *Epoll) Mod(fd int, events EpollEvent) (err error) {
	ev := &unix.EpollEvent{
		Events: uint32(events),
		Fd:     int32(fd),
	}
	ep.lock.RLock()
	defer ep.lock.RUnlock()
	if ep.closed {
		return ErrClosed
	}

	if _, ok := ep.callbacks[fd]; !ok {
		return ErrNotRegistered
	}
	return unix.EpollCtl(fd, unix.EPOLL_CTL_MOD, fd, ev)
}

const (
	maxWaitEventsBegin = 1024
	maxWaitEventsStop  = 32786
)

func (ep *Epoll) wait(OnErr func(error)) {
	defer func() {
		if err := unix.Close(ep.fd); err != nil {
			OnErr(err)
		}
		close(ep.done)
	}()
	events := make([]unix.EpollEvent, maxWaitEventsBegin)
	callbacks := make([]func(EpollEvent), maxWaitEventsBegin)

	for {
		n, err := unix.EpollWait(ep.fd, events, -1)
		if err != nil {
			if temporaryErr(err) {
				continue
			}
			OnErr(err)
			return
		}
		callbacks = callbacks[:n]
		ep.lock.RLock()
		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			if fd == ep.eventfd {
				ep.lock.RUnlock()
				return
			}
			callbacks[i] = ep.callbacks[fd]
		}
		ep.lock.RUnlock()
		for i := 0; i < n; i++ {
			if cb := callbacks[i]; cb != nil {
				cb(EpollEvent(events[i].Events))
				callbacks[i] = nil
			}
		}

		if n == maxWaitEventsBegin && n*2 < maxWaitEventsStop {
			events = make([]unix.EpollEvent, n*2)
			callbacks = make([]func(EpollEvent), n*2)
		}
	}
}

func main() {
	e, err := EpollCreate(func(e error) {

	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(e)
}
