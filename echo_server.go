package main

import (
	"goio/application"
	"goio/logger"
	"goio/msg"
	pb "goio/proto"
	"goio/protocol"

	"github.com/golang/protobuf/proto"
	"github.com/stackimpact/stackimpact-go"
)

func main() {

	agent := stackimpact.NewAgent()
	agent.Start(stackimpact.Options{
		AgentKey:       "67e7727fd85ce05c93889fccbc6f6444e045eba3",
		AppName:        "Basic Go Server",
		AppVersion:     "1.0.0",
		AppEnvironment: "production",
	})

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
				logger.Info("publish %v", msg)
				submit := &pb.Submit{}
				err := proto.Unmarshal(msg.Topic, submit)
				if err != nil {
					logger.Error("Unmarshal error %v", err)
				}
				logger.Info("Topic To %v", submit.To)

			case *protocol.MqttPingReq:
				logger.Info("ping req")
				m := &protocol.MqttPingRes{}
				m.SetChannel(msg.Channel())
				app.GetIOService().GetIOStage().Send(m)

			}
		})
		return nil
	}).Run()
}
