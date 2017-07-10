package network

import (
	"goio/logger"
	"goio/msg"
	"goio/protocol"
	"goio/util"
	"sync"
	"sync/atomic"
	//"time"
)

var (
	s *Session
	q chan msg.Message
)

func init() {
	once.Do(func() {
		//b := NewBucket(256)
		q = make(chan msg.Message, 1024)
		s = &Session{
			ctrie: &util.Map{},
			room:  &util.Map{},
			msize: 4096,
		}
		for i := uint64(0); i < s.msize; i++ {
			s.msgq[i] = make(chan *protocol.Barrage, 100)
			go procMsg(s.msgq[i])
		}

	})
}

func procMsg(q chan *protocol.Barrage) {
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
	msgq  [4096]chan *protocol.Barrage
	msize uint64
	nsize uint64
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

var l sync.Mutex

func (r *Room) handle() {
	for {
		b := <-q
		l.Lock()
		b.Channel().EncodeMessage(b)
		l.Unlock()
	}
}

func (r *Room) getMember() (uids []string, rcount int) {
	r.chs.Range(func(key, value interface{}) bool {
		uids = append(uids, key.(string))
		rcount += 1
		return true
	})
	return
}

func (r *Room) pushHost(b *protocol.Barrage) {
	q <- b
}

func (r *Room) pushMsg(barrage protocol.Barrage) {
	r.chs.Range(func(key, value interface{}) bool {
		value.(*util.Map).Range(func(key, value interface{}) bool {
			c := value.(msg.Channel)
			if c == nil || c.GetAttr("cid") == barrage.Channel().GetAttr("cid") {
				return true
			}
			b := &protocol.Barrage{}
			b.Ver = 1
			b.Op = 5
			b.Body = barrage.Body
			b.SetChannel(c)
			go c.EncodeMessage(b)
			return true
		})
		return true
	})
}

func (r *Room) addChan(cid, uid string, ch msg.Channel) {
	c, _ := r.chs.Load(uid)
	if c == nil {
		c = &util.Map{}
		r.chs.Store(uid, c)
	}
	c.(*util.Map).Store(cid, ch)
}

func (r *Room) deleteChan(cid, uid string) {
	c, _ := r.chs.Load(uid)
	if c != nil {
		ch := c.(*util.Map)
		ch.Delete(cid)
		if ch.Size() == 0 {
			r.chs.Delete(uid)
			if r.chs.Size() == 0 {
				s.room.Delete(r.rid)
			}
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
		go room.handle()
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
		r.deleteChan(cid, uid)
		r.host.Delete(cid)
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
	barrage := &protocol.Barrage{}
	barrage.Body = body
	ch := NewChannel()
	ch.SetAttr("rid", rid)
	ch.SetAttr("cid", cid)
	barrage.SetChannel(ch)
	if code == 0 {
		r.host.Range(func(key, value interface{}) bool {
			c := value.(msg.Channel)
			if c == nil || c.GetAttr("cid") == barrage.Channel().GetAttr("cid") {
				return true
			}
			b := &protocol.Barrage{}
			b.SetChannel(c)
			b.Op = 5
			b.Ver = 1
			b.Body = body
			r.pushHost(b)
			return true
		})
	} else {
		BroadcastRoom(*barrage, false)
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

func BroadcastRoom(barrage protocol.Barrage, store bool) {
	if store {
		err := util.StoreMessage("http://"+util.GetHttpConfig().Remoteaddr+"/im/"+barrage.Channel().GetAttr("rid")+"/chat", barrage.Body)

		if err != nil {
			logger.Error("util.StoreMessage error %v", err)
		}
	}
	m.Lock()
	idx := atomic.AddUint64(&s.msize, 1) % s.msize
	s.msgq[idx] <- &barrage
	m.Unlock()
}

func Close() {
	for i := uint64(0); i < s.msize; i++ {
		close(s.msgq[i])
	}
	close(q)
}
