package service

import (
	"goio/msg"
	"sync"
)

const (
	WORKER_STARTED = iota
	WORKER_STOPPED
	DISPATCHER_STARTED
	DISPATCHER_STOPPED
)

type Worker struct {
	workerPool chan chan msg.Message
	jobChannel chan msg.Message
	quit       chan struct{}
	handler    Handler
	wg         *sync.WaitGroup
	status     byte
}

type Dispatcher struct {
	queue      chan msg.Message
	workerPool chan chan msg.Message
	maxWorkers int
	workers    []*Worker
	status     byte
	stop       chan struct{}
}

func NewWorker(pool chan chan msg.Message, h Handler) *Worker {
	return &Worker{
		workerPool: pool,
		handler:    h,
		jobChannel: make(chan msg.Message),
		quit:       make(chan struct{}),
		wg:         &sync.WaitGroup{},
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	w.status = WORKER_STARTED
	go func() {
	L:
		for {
			w.workerPool <- w.jobChannel
			select {
			case msg, ok := <-w.jobChannel:
				if ok {
					w.handler.Handle(msg)
				}
			case <-w.quit:
				w.status = WORKER_STOPPED
				close(w.jobChannel)
				break L
			}
		}
		w.wg.Done()
	}()
}

func (w *Worker) Stop() {
	switch w.status {
	case WORKER_STARTED:
		close(w.quit)
	}
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func NewDispatcher(maxWorkers, maxQueue int) *Dispatcher {
	return &Dispatcher{
		queue:      make(chan msg.Message, maxQueue),
		workerPool: make(chan chan msg.Message, maxWorkers),
		maxWorkers: maxWorkers,
		workers:    make([]*Worker, maxWorkers),
		stop:       make(chan struct{}),
	}
}

func (d *Dispatcher) Run(h Handler) {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool, h)
		worker.Start()
		d.workers[i] = worker
	}
	d.status = DISPATCHER_STARTED
	go d.Dispatch()
}

func (d *Dispatcher) Stop() {
	if d.status != DISPATCHER_STOPPED {
		d.status = DISPATCHER_STOPPED
		close(d.stop)
		//close(d.queue)
		for _, worker := range d.workers {
			worker.Stop()
		}
		close(d.workerPool)
	}
}

func (d *Dispatcher) Wait() {
	for _, worker := range d.workers {
		worker.Wait()
	}
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case m := <-d.queue:
			if m != nil {
				go func(msg msg.Message) {
					jobChannel := <-d.workerPool
					if jobChannel != nil {
						jobChannel <- msg
					}
				}(m)
			}
		case <-d.stop:
			return
		}
	}
}
