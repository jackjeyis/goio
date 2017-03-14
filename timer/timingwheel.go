package timer

import (
	"errors"
	"goio/logger"
	"goio/queue"
)

type TimingNode struct {
	next   *TimingNode
	ticks  int64
	active bool
	data   interface{}
}

type TimingWheel struct {
	slots        *queue.RingBuffer
	current_slot uint64
	next_shift   uint8
	self_shift   uint8
	next         *TimingWheel
	mask         uint64
}

func NewTimingWheel(slot uint64, next *TimingWheel) *TimingWheel {
	ring, err := queue.NewRingBuffer(slot)
	if err != nil {
		logger.Error("create RingBuffer error! %v", err)
		return nil
	}
	timingWheel := &TimingWheel{
		current_slot: 0,
		slots:        ring,
		next_shift:   Log2(slot),
		self_shift:   0,
		mask:         slot - 1,
		next:         nil,
	}
	timingWheel.BindNext(next)
	return timingWheel
}

func (tw *TimingWheel) BindNext(next *TimingWheel) {
	if next == nil {
		return
	}
	tw.next = next
	next.self_shift = tw.self_shift + tw.next_shift
}

func Log2(slot uint64) (n uint8) {
	for ; slot > 0; slot >>= 2 {
		n++
	}
	return n
}

func (tw *TimingWheel) Advance(ticks int64) *TimingNode {
	if ticks == 0 {
		return nil
	}

	var (
		head      *TimingNode
		tail      *TimingNode
		slot_head *TimingNode
	)

	for ticks > 0 {
		tw.current_slot++
		if tw.slots.Index(tw.current_slot) == 0 {
			tw.PullNext()
		}
		slot_head = tw.slots.At(tw.current_slot).data.(*TimingNode)
		tw.slots.At(tw.current_slot).data = nil
		head = LinkList(&head, &tail, slot_head)
		ticks--
	}
	return head
}

func (tw *TimingWheel) AddNode(ticks uint64, node *TimingNode) error {
	var tick_slot uint64
	if ticks < tw.slots.Size() {
		tick_slot = tw.current_slot + node.ticks
		n, ok := slots.At(tick_slot).data.(*TimingNode)
		slots.At(tick_slot).data = LinkNode(n, node)
	} else {
		if tw.next != nil {
			var (
				remain_slot uint64
				offset      uint64
				loop        uint64
			)
			remain_slot = tw.slots.Size() - tw.slots.Index(tw.current_slot)
			offset = (ticks - remain_slot) * tw.mask
			loop = (ticks + slots.Index(tw.current_slot)) >> tw.next_shift
			node.ticks = (node.ticks & ^(tw.mask << tw.self_shift)) | offset<<tw.self_shift
			tw.next.AddNode(loop, node)
		} else {
			return errors.New("TimingWheel OverFlow!")
		}
	}
}

func (tw *TimingWheel) PullNext() {
	if tw.next == nil {
		return
	}
	head := tw.next.Advance(1)
	for head != nil {
		next := head.next
		var ticks uint32
		ticks = (head.ticks >> tw.self_shif) * tw.mask
		tw.AddNode(ticks, head)
		head = next
	}
}
func LinkNode(old, new *TimingNode) *TimingNode {
	new.next = old
	return new
}

func GetTail(node *TimingNode) *TimingNode {
	for node != nil && node.next != nil {
		node = node.next
	}
	return node
}

func LinkList(head, tail **TimingNode, link_node *TimingNode) *TimingNode {
	if link_node == nil {
		return *head
	}
	if *tail == nil {
		*head = link_node
		*tail = GetTail(link_node)
	} else {
		*tail.next = link_node
		*tail = GetTail(link_node)
	}
	return *head
}
