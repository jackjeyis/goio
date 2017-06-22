package service

import "goio/msg"
import "sync"

type Handler interface {
	Handle(msg.Message)
}

type ServiceHandler struct {
}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

func (sh *ServiceHandler) Handle(msg msg.Message) {
	msg.Channel().Serve(msg)
}

type IOHandler struct {
}

func NewIOHandler() *IOHandler {
	return &IOHandler{}
}
var m sync.Mutex
func (ih *IOHandler) Handle(msg msg.Message) {
	m.Lock()
	msg.Channel().EncodeMessage(msg)
	msg.Channel().OnWrite()
	m.Unlock()
}
