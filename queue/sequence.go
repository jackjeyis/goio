package queue

import (
	"errors"
	"runtime"
	"sync/atomic"
)

type node struct {
	pos  uint64
	data interface{}
}

type Sequence struct {
	padding0       [8]uint64
	produce        uint64
	padding1       [8]uint64
	consume        uint64
	padding2       [8]uint64
	mask, disposed uint64
	padding3       [8]uint64
	nodes          []*node
}

func roundUp(num uint64) uint64 {
	num--
	num |= num >> 1
	num |= num >> 2
	num |= num >> 4
	num |= num >> 8
	num |= num >> 16
	num |= num >> 32
	num++
	return num
}

func (s *Sequence) init(size uint64) {
	size = roundUp(size)
	s.nodes = make([]*node, size)
	for i := uint64(0); i < size; i++ {
		s.nodes[i] = &node{pos: i}
	}
	s.mask = size - 1
}

func NewSequence(size uint64) *Sequence {
	seq := &Sequence{}
	seq.init(size)
	return seq
}

func (s *Sequence) put(item interface{}, offer bool) (bool, error) {
	var (
		n *node
		i int
	)
	pos := atomic.LoadUint64(&s.produce)
L:
	for {
		if s.Disposed() {
			return false, errors.New("disposed")
		}
		n = s.nodes[pos&s.mask]
		seq := atomic.LoadUint64(&n.pos)

		switch diff := seq - pos; {
		case diff == 0:
			if atomic.CompareAndSwapUint64(&s.produce, pos, pos+1) {
				break L
			}
		case diff < 0:
			return false, errors.New("Sequence in compromised state during a get operation.")
		default:
			pos = atomic.LoadUint64(&s.produce)
		}

		if offer {
			return false, nil
		}

		if i == 10000 {
			runtime.Gosched()
			i = 0
		} else {
			i++
		}
	}
	n.data = item
	atomic.StoreUint64(&n.pos, pos+1)
	return true, nil
}

func (s *Sequence) Put(item interface{}) error {
	_, err := s.put(item, false)
	return err
}

func (s *Sequence) Offer(item interface{}) error {
	_, err := s.put(item, true)
	return err
}

func (s *Sequence) Get() (interface{}, error) {
	var n *node
	pos := atomic.LoadUint64(&s.consume)
	//i := 0
L:
	for {
		if s.Disposed() {
			return nil, errors.New("disposed")
		}
		n = s.nodes[pos&s.mask]
		seq := atomic.LoadUint64(&n.pos)

		switch diff := seq - (pos + 1); {
		case diff == 0:
			if atomic.CompareAndSwapUint64(&s.consume, pos, pos+1) {
				break L
			}
		case diff < 0:
			return nil, errors.New("Sequence in compromised state during a get operation.")
		case diff > 0:
			pos = atomic.LoadUint64(&s.consume)
		}

		/*if i == 10000 {
			runtime.Gosched()
			i = 0
		} else {
			i++
		}*/
		return nil, errors.New("no item")
	}

	data := n.data
	n.data = nil
	atomic.StoreUint64(&n.pos, pos+s.mask+1)
	return data, nil
}

func (s *Sequence) Len() uint64 {
	return atomic.LoadUint64(&s.produce) - atomic.LoadUint64(&s.consume)
}

func (s *Sequence) Cap() uint64 {
	return uint64(len(s.nodes))
}

func (s *Sequence) Disposed() bool {
	return atomic.CompareAndSwapUint64(&s.disposed, 1, 0)
}
