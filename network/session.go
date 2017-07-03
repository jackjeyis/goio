package network

import (
	"goio/logger"
	"goio/msg"
	"goio/protocol"
	"goio/util"
	"sync"
	"sync/atomic"
)

var (
	s *Session
)

func init() {
	once.Do(func() {
		//b := NewBucket(256)
		s = &Session{
			ctrie: &util.Map{},
			room:  &util.Map{},
			size:  1024,
		}
		for i := uint64(0); i < s.size; i++ {
			s.queue[i] = make(chan *protocol.Barrage, s.size/2)
			go process(s.queue[i])
		}
	})
}

func process(q chan *protocol.Barrage) {
	var b *protocol.Barrage
	for {
		b = <-q
		if b != nil {
			PushRoom(*b)
		}
	}
}

type Bucket struct {
	buckets []*Session
	size    int
}

func NewBucket(size int) *Bucket {
	b := &Bucket{
		buckets: make([]*Session, size),
		size:    size,
	}

	for i := 0; i < size; i++ {
		b.buckets[i] = &Session{ctrie: &util.Map{}, room: &util.Map{}}
	}

	return b
}

func (b *Bucket) HashSession(key string) *Session {
	idx := 0 /*cityhash.CityHash32([]byte(key), uint32(len(key))) % b.size*/
	return b.buckets[idx]
}

type Session struct {
	ctrie *util.Map
	room  *util.Map
	queue [1024]chan *protocol.Barrage
	size  uint64
}

type Room struct {
	chs  *util.Map
	rid  string
	host *util.Map
}

func NewRoom(id string) *Room {
	return &Room{
		chs:  &util.Map{},
		rid:  id,
		host: &util.Map{},
	}
}

func (r *Room) addChan(cid, uid string, ch msg.Channel) {
	c, _ := r.chs.Load(uid)
	if c == nil {
		c = &util.Map{}
		r.chs.Store(uid, c)
	}
	c.(*util.Map).Store(cid, ch)
}

func (r *Room) getMember() (uids []string, rcount int) {
	r.chs.Range(func(key, value interface{}) bool {
		uids = append(uids, key.(string))
		rcount += 1
		return true
	})
	return
}

/*func (r *Room) getMember() (chans []msg.Channel, uids []string, rcount int) {
	for entry := range r.chs.Iterator(nil) {
		for item := range entry.Value.(*util.Ctrie).Iterator(nil) {
			chans = append(chans, item.Value.(msg.Channel))
		}
		uids = append(uids, string(entry.Key))
		rcount += 1
	}
	return
}*/

func (r *Room) pushMsg(barrage protocol.Barrage) {
	r.chs.Range(func(key, value interface{}) bool {
		value.(*util.Map).Range(func(key, value interface{}) bool {
			c := value.(msg.Channel)
			if c == nil || c.GetAttr("cid") == barrage.Channel().GetAttr("cid") {
				return true
			}

			barrage.Ver = 1
			barrage.Op = 5
			barrage.SetChannel(c)
			c.EncodeMessage(&barrage)
			return true
		})
		return true
	})
}

/*func (r *Room) pushMsg(barrage protocol.Barrage) {
	for entry := range r.chs.Iterator(nil) {
		for item := range entry.Value.(*util.Ctrie).Iterator(nil) {
			c := item.Value.(msg.Channel)
			if c == nil || c.GetAttr("cid") == barrage.Channel().GetAttr("cid") {
				continue
			}
			barrage.Ver = 1
			barrage.Op = 5
			barrage.SetChannel(c)
			c.EncodeMessage(&barrage)
		}
	}
}*/

func (r *Room) deleteChan(cid, uid string) {
	c, _ := r.chs.Load(uid)
	if c != nil {
		ch := c.(*util.Map)
		ch.Delete(cid)
		if ch.Size() == 0 {
			r.chs.Delete(uid)
		}
	}
}

func (s *Session) insert(cid string, ch msg.Channel, host int) {
	var room *Room
	s.ctrie.Store(cid, ch)
	rid := ch.GetAttr("rid")
	r, _ := s.room.Load(rid)
	if r == nil {
		room = NewRoom(rid)
		s.room.Store(rid, room)
	} else {
		room = r.(*Room)
	}
	if host != 2 {
		room.host.Store(cid, ch)
	}
	room.addChan(cid, ch.GetAttr("uid"), ch)
}

func (s *Session) delete(cid, uid, rid string) {
	s.ctrie.Delete(cid)
	r := s.roomer(rid)
	if r != nil {
		if r.chs.Size() == 0 {
			s.room.Delete(rid)
		} else {
			r.deleteChan(cid, uid)
			r.host.Delete(cid)
		}
	}
}

func (s *Session) find(cid string) msg.Channel {
	ch, ok := s.ctrie.Load(cid)
	if ok {
		return ch.(msg.Channel)
	}
	return nil
}

func (s *Session) roomer(rid string) *Room {
	r, ok := s.room.Load(rid)
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

/*func GetRoomSession(rid string) []msg.Channel {
	r := s.roomer(rid)
	if r == nil {
		return nil
	}
	chans, _, _ := r.getMember()
	return chans
}*/

func GetRoomStatus(rid string) (int, []string) {
	r := s.roomer(rid)
	if r == nil {
		return 0, nil
	}
	uids, count := r.getMember()
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
	notify.Uid = uid
	notify.Rid = rid
	notify.Code = code
	notify.Time = util.GetMillis()
	body, err := util.EncodeJson(notify)
	err = util.StoreMessage("http://"+util.GetHttpConfig().Remoteaddr+"/im/"+rid+"/view_record", body)

	if err != nil {
		logger.Error("util.StoreMessage error %v", err)
	}
	barrage := protocol.Barrage{}
	barrage.Body = body
	ch := NewChannel()
	ch.SetAttr("rid", rid)
	ch.SetAttr("cid", cid)
	barrage.SetChannel(ch)
	if code == 0 {
		r.host.Range(func(key, value interface{}) bool {
			ch := value.(msg.Channel)
			barrage.SetChannel(ch)
			ch.EncodeMessage(&barrage)
			return true
		})
	} else {
		BroadcastRoom(barrage, false)
	}
}

func PushRoom(barrage protocol.Barrage) {
	r := s.roomer(barrage.Channel().GetAttr("rid"))
	if r == nil {
		return
	}
	r.pushMsg(barrage)
}

func Push(barrage protocol.Barrage) {
	barrage.Channel().GetIOService().Serve(&barrage)
}

var m sync.Mutex

func BroadcastRoom(barrage protocol.Barrage, store bool) {
	if store {
		err := util.StoreMessage("http://"+util.GetHttpConfig().Remoteaddr+"/im/"+barrage.Channel().GetAttr("rid")+"/chat", barrage.Body)

		if err != nil {
			logger.Error("util.StoreMessage error %v", err)
		}
	}
	m.Lock()
	idx := atomic.AddUint64(&s.size, 1) % s.size
	s.queue[idx] <- &barrage
	m.Unlock()
}

func Close() {
	for i := uint64(0); i < s.size; i++ {
		close(s.queue[i])
	}
}
