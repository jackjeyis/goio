package main

import (
	"goio/application"
	"goio/logger"
	"goio/msg"
	"goio/protocol"
	//"github.com/stackimpact/stackimpact-go"
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

	app := &application.GenericApplication{}
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
			barrage, ok := msg.(*protocol.Barrage)
			if !ok {
				return
			}
			if msg.Channel().GetAttr("status") != "OK" {
				if barrage.Op != 7 {
					logger.Info("handshake fail %v", barrage)
					msg.Channel().Close()
					return
				}
				logger.Info("auth %v", string(barrage.Body))
				msg.Channel().SetAttr("status", "OK")
				barrage.Body = nil
				barrage.Op = 8
				app.GetIOService().GetIOStage().Send(barrage)
			}
			switch barrage.Op {
			case 2:
				barrage.Op = 3
				app.GetIOService().GetIOStage().Send(barrage)
			case 4:
				barrage.Op = 5
				app.GetIOService().GetIOStage().Send(barrage)
			}
		})
		return nil
	}).Run()

}
