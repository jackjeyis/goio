package network

import (
	"goio/logger"
	"goio/msg"
	"goio/protocol"
	"goio/util"
	"strconv"
)

var (
	s *Session
)

func init() {
	once.Do(func() {
		s = &Session{
			ctrie: util.New(nil),
			room:  util.New(nil),
		}
	})
}

type Bucket struct {
	buckets []*Session
}
type Session struct {
	ctrie *util.Ctrie
	room  *util.Ctrie
}

type Room struct {
	chs  *util.Ctrie
	rid  string
	host *util.Ctrie
}

func NewRoom(id string) *Room {
	return &Room{
		chs:  util.New(nil),
		rid:  id,
		host: util.New(nil),
	}
}

func (r *Room) addChan(cid, uid string, ch msg.Channel) {
	c, _ := r.chs.Lookup([]byte(uid))
	if c == nil {
		c = util.New(nil)
		r.chs.Insert([]byte(uid), c)
	}
	c.(*util.Ctrie).Insert([]byte(cid), ch)
}

func (r *Room) getMember() (chans []msg.Channel, uids []string, rcount int) {
	for entry := range r.chs.Iterator(nil) {
		for item := range entry.Value.(*util.Ctrie).Iterator(nil) {
			chans = append(chans, item.Value.(msg.Channel))
		}
		uids = append(uids, string(entry.Key))
		rcount += 1
	}
	return
}

func (r *Room) deleteChan(cid, uid string) {
	c, _ := r.chs.Lookup([]byte(uid))
	if c != nil {
		ch := c.(*util.Ctrie)
		ch.Remove([]byte(cid))
		if ch.Size() == 0 {
			r.chs.Remove([]byte(uid))
		}
	}
}

func (s *Session) insert(cid string, ch msg.Channel, host int) {
	var room *Room
	s.ctrie.Insert([]byte(cid), ch)
	rid := ch.GetAttr("rid")
	r, _ := s.room.Lookup([]byte(rid))
	if r == nil {
		room = NewRoom(rid)
		s.room.Insert([]byte(rid), room)
	} else {
		room = r.(*Room)
	}
	if host != 2 {
		room.host.Insert([]byte(cid), ch)
	}
	room.addChan(cid, ch.GetAttr("uid"), ch)
}

func (s *Session) delete(cid, uid, rid string) {
	ch, _ := s.ctrie.Remove([]byte(cid))
	if ch == nil {
		return
	}
	r := s.roomer(rid)
	if r != nil {
		if r.chs.Size() == 0 {
			s.room.Remove([]byte(rid))
		} else {
			r.deleteChan(cid, uid)
			r.host.Remove([]byte(cid))
		}
	}
}

func (s *Session) find(cid string) msg.Channel {
	ch, ok := s.ctrie.Lookup([]byte(cid))
	if ok {
		return ch.(msg.Channel)
	}
	return nil
}

func (s *Session) roomer(rid string) *Room {
	r, ok := s.room.Lookup([]byte(rid))
	if !ok {
		return nil
	}
	return r.(*Room)
}

func Register(cid string, ch msg.Channel, host int) {
	s.insert(cid, ch, host)
}

func UnRegister(cid, uid, rid string) {
	s.delete(cid, uid, rid)
}

func GetSession(cid string) msg.Channel {
	return s.find(cid)
}

func GetRoomSession(rid string) []msg.Channel {
	r := s.roomer(rid)
	if r == nil {
		return nil
	}
	chans, _, _ := r.getMember()
	return chans
}

func GetRoomStatus(rid string) (int, []string) {
	r := s.roomer(rid)
	if r == nil {
		return 0, nil
	}
	_, uids, count := r.getMember()
	return count, uids
}

func NotifyHost(rid, cid, uid string, code int8) {
	r := s.roomer(rid)
	if r == nil {
		return
	}
	notify := protocol.Notify{}
	notify.Id = util.UUID()
	notify.Ct = 90010
	Uid, _ := strconv.ParseInt(uid, 10, 64)
	notify.Uid = Uid
	Rid, _ := strconv.ParseInt(rid, 10, 32)
	notify.Rid = int32(Rid)
	notify.Code = code
	notify.Time = util.GetMillis()
	body, err := util.EncodeJson(notify)
	err = util.StoreMessage("http://"+util.GetHttpConfig().Remoteaddr+"/im/"+rid+"/view_record", body)

	if err != nil {
		logger.Error("util.StoreMessage error %v", err)
	}
	if code == 0 {
		for item := range r.host.Iterator(nil) {
			Push(item.Value.(msg.Channel), body)
		}
	} else {
		PushRoom(rid, cid, body)
	}
}

func PushRoom(rid, cid string, body []byte) {
	chans := GetRoomSession(rid)
	if chans == nil {
		return
	}
	for _, c := range chans {
		if c == nil || c.GetAttr("cid") == cid {
			continue
		}
		Push(c, body)
	}
}

func Push(ch msg.Channel, body []byte) {
	msg := &protocol.Barrage{}
	msg.Op = 5
	msg.Ver = 1
	msg.Body = body
	msg.SetChannel(ch)
	ch.GetIOService().Serve(msg)
}

func BroadcastRoom(rid, cid string, body []byte, store bool) {
	if store {
		err := util.StoreMessage("http://"+util.GetHttpConfig().Remoteaddr+"/im/"+rid+"/chat", body)

		if err != nil {
			logger.Error("util.StoreMessage error %v", err)
		}
	}
	PushRoom(rid, cid, body)
}
