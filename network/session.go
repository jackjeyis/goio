package network

import (
	"goio/msg"
	"goio/util"
)

var (
	s *Session
)

func init() {
	once.Do(func() {
		s = &Session{
			ctrie: util.New(nil),
		}
	})
}

type Session struct {
	ctrie *util.Ctrie
}

func (s *Session) insert(cid string, ch msg.Channel) {
	s.ctrie.Insert([]byte(cid), ch)
}

func (s *Session) delete(cid string) {
	s.ctrie.Remove([]byte(cid))
}

func (s *Session) find(cid string) msg.Channel {
	ch, ok := s.ctrie.Lookup([]byte(cid))
	if ok {
		return ch.(msg.Channel)
	}
	return nil
}

func Register(cid string, ch msg.Channel) {
	s.insert(cid, ch)
}

func UnRegister(cid string) {
	s.delete(cid)
}

func GetSession(cid string) msg.Channel {
	return s.find(cid)
}
