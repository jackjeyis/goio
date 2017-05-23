package network

import (
	"goio/hp"
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
	cid_uid *util.Ctrie
	rid     string
	host    msg.Channel
	admin   msg.Channel
}

func NewRoom(id string) *Room {
	return &Room{
		cid_uid: util.New(nil),
		rid:     id,
	}
}

func (r *Room) addUid(cid string, uid int64) {
	r.cid_uid.Insert([]byte(cid), uid)
}

func (r *Room) getMember() (chans []msg.Channel, uids []int64, rcount int) {
	for entry := range r.cid_uid.Iterator(nil) {
		chans = append(chans, GetSession(string(entry.Key)))
		uids = append(uids, entry.Value.(int64))
		rcount += 1
	}
	return
}

func (r *Room) deleteCid(cid string) {
	r.cid_uid.Remove([]byte(cid))
}

func (s *Session) insert(cid string, ch msg.Channel, host int) {
	var room *Room
	s.ctrie.Insert([]byte(cid), ch)
	rid := ch.GetAttr("rid").(string)
	r, _ := s.room.Lookup([]byte(rid))
	if r == nil {
		room = NewRoom(rid)
		s.room.Insert([]byte(rid), room)
	} else {
		room = r.(*Room)
	}
	if host == 0 {
		room.admin = ch
	} else if host == 1 {
		room.host = ch
	}
	room.addUid(cid, ch.GetAttr("uid").(int64))
}

func (s *Session) delete(cid, rid string) {
	ch, _ := s.ctrie.Remove([]byte(cid))
	if ch == nil {
		return
	}
	r := s.roomer(rid)
	if r != nil {
		if r.cid_uid.Size() == 0 {
			s.room.Remove([]byte(rid))
		} else {
			r.deleteCid(ch.(msg.Channel).GetAttr("cid").(string))
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

func UnRegister(cid, rid string) {
	s.delete(cid, rid)
}

func IsRegister(cid string) bool {
	return GetSession(cid) != nil
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

func GetRoomStatus(rid string) (int, []int64) {
	r := s.roomer(rid)
	if r == nil {
		return 0, nil
	}
	_, uids, count := r.getMember()
	return count, uids
}

func NotifyHost(rid string, code int8) {
	r := s.roomer(rid)
	if r == nil {
		return
	}

	if r.host != nil {
		Push(rid, r.host, code)
	}

	if r.admin != nil {
		Push(rid, r.admin, code)
	}
}

func Push(rid string, ch msg.Channel, code int8) {
	notify := protocol.Notify{}
	notify.Id = util.UUID()
	notify.Ct = 90010
	notify.Uid = ch.GetAttr("uid").(int64)
	Rid, _ := strconv.ParseInt(rid, 10, 32)
	notify.Rid = int32(Rid)
	notify.Code = code
	body, err := hp.EncodeJson(notify)
	err = hp.StoreMessage(rid, body)

	if err != nil {
		logger.Error("hp.StoreMessage error %v", err)
	}

	msg := &protocol.Barrage{}
	msg.Op = 5
	msg.Ver = 1
	msg.Body = body
	msg.SetChannel(ch)
	ch.GetIOService().Serve(msg)
}

func BroadcastRoom(rid, cid string, body []byte) {
	chans := GetRoomSession(rid)
	if chans == nil {
		return
	}
	err := hp.StoreMessage(rid, body)

	if err != nil {
		logger.Error("hp.StoreMessage error %v", err)
	}
	for _, c := range chans {
		if c == nil || c.GetAttr("cid").(string) == cid {
			continue
		}
		msg := &protocol.Barrage{}
		msg.Op = 5
		msg.Ver = 1
		msg.Body = body
		msg.SetChannel(c)
		c.GetIOService().Serve(msg)
	}
}
