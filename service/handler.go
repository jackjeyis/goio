package service

import "goio/msg"

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

func (ih *IOHandler) Handle(msg msg.Message) {
	msg.Channel().EncodeMessage(msg)
}
