package service

import (
	"goio/logger"
	"goio/msg"
)

type Handler interface {
	Handle(msg.Message)
}

type ServiceHandler struct {
}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{}
}

func (sh *ServiceHandler) Handle(msg msg.Message) {
	logger.Info("Service Handle")
	msg.ProtoMsg().Serve(msg)
}

type IOHandler struct {
}

func NewIOHandler() *IOHandler {
	return &IOHandler{}
}

func (ih *IOHandler) Handle(msg msg.Message) {
	logger.Info("IO handle")
	msg.ProtoMsg().EncodeMessage(msg)
}
