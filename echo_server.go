package main

import (
	"encoding/json"
	"goio/application"
	"goio/hp"
	"goio/logger"
	"goio/msg"
	"goio/network"
	"goio/protocol"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
	//"github.com/stackimpact/stackimpact-go"
)

type Auth struct {
	Rid   int64
	Token string
	Cid   string
}

type AuthReply struct {
	Code int
	Msg  string
}

var (
	app  *application.GenericApplication
	auth Auth
	uid  int64
)

type TermType byte

const (
	Android = TermType(iota)
	IOS
	PC
)

func main() {

	/*agent := stackimpact.NewAgent()
	agent.Start(stackimpact.Options{
		AgentKey:       "67e7727fd85ce05c93889fccbc6f6444e045eba3",
		AppName:        "Basic Go Server",
		AppVersion:     "1.0.0",
		AppEnvironment: "production",
	})
	*/

	app = &application.GenericApplication{}
	app.SetOnStart(func() error {
		app.RegisterService("mqtt_handler", func(msg msg.Message) {
			switch msg := msg.(type) {
			case *protocol.MqttConnect:
				logger.Info("connect success!")
				//network.Register(msg.ClientId, msg.Channel())
				m := &protocol.MqttConnAck{RetCode: protocol.RetCodeAccepted}
				//m.SetChannel(network.GetSession(msg.ClientId))
				m.SetChannel(msg.Channel())
				app.GetIOService().GetIOStage().Send(m)
			case *protocol.MqttPublish:
				//logger.Info("publish %v", msg)
				//m := &protocol.MqttPubAck{MsgId: msg.MsgId}
				//m.SetChannel(msg.Channel())
				//app.GetIOService().GetIOStage().Send(m)
				/*submit := &pb.Submit{}
				err := proto.Unmarshal(msg.Topic, submit)
				if err != nil {
					logger.Error("Unmarshal error %v", err)
				}
				logger.Info("Topic To %v", submit.To)
				*/
				app.GetIOService().GetIOStage().Send(msg)

			case *protocol.MqttPingReq:
				logger.Info("ping req")
				m := &protocol.MqttPingRes{}
				m.SetChannel(msg.Channel())
				app.GetIOService().GetIOStage().Send(m)

			}
		})
		app.RegisterService("barrage_handler", func(msg msg.Message) {
			barrage := msg.(*protocol.Barrage)
			if barrage == nil {
				return
			}

			if barrage.Channel().GetAttr("status") != "OK" {
				if barrage.Op != 7 {
					logger.Info("handshake fail!")
					barrage.Channel().Close()
					return
				}
				logger.Info("auth %v", string(barrage.Body))
				if err := json.Unmarshal(barrage.Body, &auth); err != nil {
					logger.Error("json.Unmarshal error %v,body %v", err, string(barrage.Body))
					barrage.Channel().Close()
					return
				}
				barrage.Op = 8
				rid := strconv.FormatInt(auth.Rid, 10)
				res, err := hp.Get("http://172.16.6.135:8998/im/" + rid + "/_check?token=" + auth.Token)
				if err != nil {
					logger.Error("check token api error %v", err)
					barrage.Channel().Close()
					return
				}
				var reply AuthReply

				switch res.Code {
				case 0:
					reply.Code = 0
					reply.Msg = "鉴权成功!"
				case 403:
					reply.Code = 1
					reply.Msg = "token 校验失败!"
				case 1001:
					reply.Code = 2
					reply.Msg = "该直播已被管理员删除!"
				case 1002:
					reply.Code = 3
					reply.Msg = "您已被管理员移出，无权限观看!"
				}

				if network.IsRegister(auth.Cid) {
					reply.Code = 4
					reply.Msg = "该用户已连接!"
				}

				body, _ := hp.EncodeJson(reply)
				barrage.Body = body
				if reply.Code != 0 {
					barrage.Channel().GetIOService().Serve(barrage)
					time.Sleep(1 * time.Second)
					barrage.Channel().Close()
					return
				}
				barrage.Channel().SetAttr("status", "OK")
				barrage.Channel().SetAttr("cid", auth.Cid)
				barrage.Channel().SetAttr("uid", int64(res.Data.UserId))
				barrage.Channel().SetAttr("rid", rid)
				barrage.Channel().SetAttr("ct", IOS)
				network.Register(auth.Cid, barrage.Channel(), res.Data.Role)
				barrage.Channel().GetIOService().Serve(barrage)
				network.NotifyHost(rid, 1)
			}
			switch barrage.Op {
			case 2:
				barrage.Op = 3
				barrage.Channel().SetDeadline(6)
				barrage.Channel().GetIOService().Serve(barrage)
			case 4:
				network.BroadcastRoom(barrage.Channel().GetAttr("rid").(string), barrage.Channel().GetAttr("cid").(string), barrage.Body)
			}
		})

		go func() {
			httpServeMux := http.NewServeMux()
			httpServeMux.HandleFunc("/1/push/room", PushRoom)
			httpServeMux.HandleFunc("/1/get/room", GetRoom)
			httpServer := &http.Server{
				Handler:      httpServeMux,
				ReadTimeout:  time.Duration(5) * time.Second,
				WriteTimeout: time.Duration(5) * time.Second,
			}
			httpServer.SetKeepAlivesEnabled(true)
			addr := app.GetConfigManager().GetHttpConfig().Addr
			ln, err := net.Listen("tcp", "0.0.0.0:7172")
			if err != nil {
				logger.Error("net.Listen %s error %v", addr)
				panic(err)
			}
			logger.Info("start http listen %s", addr)
			if err = httpServer.Serve(ln); err != nil {
				logger.Error("httpServer.Serve() error %v", err)
				panic(err)
			}
		}()
		return nil
	}).Run()

}

func retWrite(w http.ResponseWriter, r *http.Request, res map[string]interface{}, start time.Time) {
	data, err := json.Marshal(res)
	if err != nil {
		logger.Error("json.Marshal error %v", err)
		return
	}
	if _, err := w.Write(data); err != nil {
		logger.Error("w.Write error %v", err)
	}
	logger.Info("req: \"%s\",get: ip:\"%s\",time:\"%fs\"", r.URL.String(), r.RemoteAddr, time.Now().Sub(start).Seconds())
}

func PushRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	var (
		bodyBytes []byte
		err       error
		param     = r.URL.Query()
		res       = map[string]interface{}{"ret": 1}
	)

	defer retWrite(w, r, res, time.Now())
	w.Header().Set("Content-Type", "application/json")
	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		logger.Error("ioutil.ReadAll failed %v", err)
		res["ret"] = 65535
		return
	}
	network.BroadcastRoom(param.Get("rid"), "", bodyBytes)
	return
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	var (
		param = r.URL.Query()
		res   = map[string]interface{}{"ret": 1}
	)

	defer retWrite(w, r, res, time.Now())
	w.Header().Set("Content-Type", "application/json")
	logger.Info("rid %v", param.Get("rid"))
	count, uids := network.GetRoomStatus(param.Get("rid"))
	if uids == nil {
		res["ret"] = 2
		return
	}
	res["data"] = uids
	res["total"] = count
	return
}
