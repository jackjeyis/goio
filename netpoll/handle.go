package netpoll

import (
	"net"
	"os"
	"syscall"
)

type filer interface {
	File() (*os.File, error)
}

type Desc struct {
	file  *os.File
	event Event
}

func NewDesc(fd uintptr, ev Event) *Desc {
	return &Desc{os.NewFile(fd, ""), ev}
}

func (desc *Desc) Close() error {
	return desc.file.Close()
}

func (desc *Desc) fd() int {
	return int(desc.file.Fd())
}

func Must(desc *Desc, err error) *Desc {
	if err != nil {
		panic(err)
	}
	return desc
}

func HandleRead(conn net.Conn) (*Desc, error) {
	return Handle(conn, EventRead|EventEdgeTriggered)
}

func HandleReadOnce(conn net.Conn) (*Desc, error) {
	return Handle(conn, EventRead|EventOneShot)
}

func HandleWrite(conn net.Conn) (*Desc, error) {
	return Handle(conn, EventWrite|EventEdgeTriggered)
}

func HandleWriteOnce(conn net.Conn) (*Desc, error) {
	return Handle(conn, EventWrite|EventOneShot)
}

func HandleReadWrite(conn net.Conn) (*Desc, error) {
	return Handle(conn, EventRead|EventWrite|EventEdgeTriggered)
}

func Handle(conn net.Conn, event Event) (*Desc, error) {
	desc, err := handle(conn, event)
	if err != nil {
		return nil, err
	}
	if err = syscall.SetNonblock(desc.fd(), true); err != nil {
		return nil, os.NewSyscallError("setnonblock", err)
	}
	return desc, nil
}

func HandleListener(ln net.Listener, event Event) (*Desc, error) {
	return handle(ln, event)
}

func handle(x interface{}, event Event) (*Desc, error) {
	f, ok := x.(filer)
	if !ok {
		return nil, ErrNotFd
	}

	file, err := f.File()
	if err != nil {
		return nil, err
	}

	return &Desc{
		file:  file,
		event: event,
	}, nil
}
