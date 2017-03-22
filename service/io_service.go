package service

import (
	"goio/msg"
	"os"
	"os/signal"
	"syscall"
)

type Stage interface {
	Start()
	Wait()
	Stop()
	Send(msg.Message)
}

type ServiceStage struct {
	stage        *Dispatcher
	srv_handle   *ServiceHandler
	worker_count int
	queue_size   int
}

func (ss *ServiceStage) Start() {
	ss.stage = NewDispatcher(ss.worker_count, ss.queue_size)
	ss.srv_handle = NewServiceHandler()
	ss.stage.Run(ss.srv_handle)
}

func (ss *ServiceStage) Wait() {
	ss.stage.Wait()
}

func (ss *ServiceStage) Stop() {
	ss.stage.Stop()
}

func (ss *ServiceStage) Send(msg msg.Message) {
	if ss.stage.status == DISPATCHER_STARTED {
		ss.stage.queue <- msg
	}
}

type IOStage struct {
	stage        *Dispatcher
	io_handle    *IOHandler
	worker_count int
	queue_size   int
}

func (io *IOStage) Start() {
	io.stage = NewDispatcher(io.worker_count, io.queue_size)
	io.io_handle = NewIOHandler()
	io.stage.Run(io.io_handle)
}

func (io *IOStage) Wait() {
	io.stage.Wait()
}

func (io *IOStage) Stop() {
	io.stage.Stop()
}

func (io *IOStage) Send(msg msg.Message) {
	if io.stage.status == DISPATCHER_STARTED {
		io.stage.queue <- msg
	}
}

type IOServiceConfig struct {
	Srvworker int
	Srvqueue  int

	Ioworker int
	Ioqueue  int
}

type IOService struct {
	service_stage *ServiceStage
	io_stage      *IOStage
}

func (s *IOService) Init(c IOServiceConfig) error {
	s.service_stage = &ServiceStage{worker_count: c.Srvworker, queue_size: c.Srvqueue}
	s.io_stage = &IOStage{worker_count: c.Ioworker, queue_size: c.Ioqueue}
	return nil
}

func (s *IOService) Start() {
	s.service_stage.Start()
	s.io_stage.Start()
}

func (s *IOService) Run() {
	s.service_stage.Wait()
}

func (s *IOService) CleanUp() {
	s.io_stage.Stop()
	s.io_stage.Wait()
}

func (s *IOService) Stop() {
	s.service_stage.Stop()
}

func (s *IOService) GetServiceStage() Stage {
	return s.service_stage
}

func (s *IOService) GetIOStage() Stage {
	return s.io_stage
}

func (s *IOService) HandleSignal() {
	go func() {
		for {
			sigs := make(chan os.Signal)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			for sig := range sigs {
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					s.Stop()
				}
			}
		}
	}()
}
