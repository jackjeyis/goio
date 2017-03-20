package main

import (
	"goio/application"
	"goio/logger"
	"goio/msg"
	"goio/network"
	pb "goio/proto"
	"goio/protocol"

	"github.com/golang/protobuf/proto"
)

func main() {
	app := &application.GenericApplication{}
	app.SetOnStart(func() error {
		app.RegisterService("mqtt_handler", func(msg msg.Message) {
			switch msg := msg.(type) {
			case *protocol.MqttConnect:
				logger.Info("connect success!")
				network.Register(msg.ClientId, msg.Channel())
				m := &protocol.MqttConnAck{}
				m.SetChannel(network.GetSession(msg.ClientId))
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
