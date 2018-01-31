package netpoll

import (
	"errors"
	"log"
)

var (
	ErrNotFd         = errors.New("could not get file discriptor!")
	ErrClosed        = errors.New("poller instance is closed!")
	ErrRegistered    = errors.New("file discriptor is already registered in poller instance")
	ErrNotRegistered = errors.New("file discriptor not registered in poller instance")
)

type Event uint16

const (
	EventRead          Event = 0x1
	EventWrite               = 0x2
	EventOneShot             = 0x4
	EventEdgeTriggered       = 0x8
	EventHup                 = 0x10
	EventReadHup             = 0x20
	EventWriteHup            = 0x40
	EventErr                 = 0x80
	EventPollerClosed        = 0x8000
)

func (ev Event) String() (str string) {
	name := func(event Event, name string) {
		if ev&event == 0 {
			return
		}
		if str != "" {
			str += "|"
		}
		str += name
	}

	name(EventRead, "EventRead")
	name(EventWrite, "EventWrite")
	name(EventOneShot, "EventOneShot")
	name(EventEdgeTriggered, "EventEdgeTriggered")
	name(EventReadHup, "EventReadHup")
	name(EventWriteHup, "EventWriteHup")
	name(EventHup, "EventHup")
	name(EventErr, "EventErr")
	name(EventPollerClosed, "EventPollerClosed")

	return
}

type Poller interface {
	Start(*Desc, CallbackFn) error
	Stop(*Desc) error
	Resume(*Desc) error
}

//type Callback func(Event)

// CallbackFn is a function that will be called on kernel i/o event
// notification.
type CallbackFn func(Event)

// Config contains options for Poller configuration.
type Config struct {
	// OnWaitError will be called from goroutine, waiting for events.
	OnWaitError func(error)
}

func (c *Config) withDefaults() (config Config) {
	if c != nil {
		config = *c
	}
	if config.OnWaitError == nil {
		config.OnWaitError = defaultOnWaitError
	}
	return config
}

func defaultOnWaitError(err error) {
	log.Printf("netpoll: wait loop error: %s", err)
}
