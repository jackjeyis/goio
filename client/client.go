package client

import (
	"goio/logger"
	"goio/proto"
	"goio/protocol"
	"goio/queue"
	"net"
	"time"

	pb "github.com/golang/protobuf/proto"
)

type Client struct {
	conn    net.Conn
	timeout time.Duration
	err     error
	buf     *queue.IOBuffer
}

func NewClient() *Client {
	return &Client{
		timeout: 2 * time.Second,
		buf:     queue.NewIOBuffer(false),
	}
}

func (c *Client) Connect(addr string) {
	c.conn, c.err = net.DialTimeout("tcp", addr, c.timeout)
	if c.err != nil {
		logger.Error("net.Dial")
	}
	ping := &protocol.PingProtocol{}
	msg := &protocol.PingPackage{}
	header := protocol.PingHeader{
		Op:  1,
		Ver: 0,
	}

	connect := &proto.Connect{
		KeepAlive: 2,
		Token:     []byte{'a'},
		LoginType: false,
		ClientID:  []byte{'i'},
	}
	data, err := pb.Marshal(connect)
	if err != nil {
		logger.Error("pb.Marshal error %v", err)
	}
	msg.PingHeader = header
	msg.Body = data
	ping.Encode(msg, c.buf)
	c.conn.Write(c.buf.Buffer()[c.buf.GetRead():c.buf.GetWrite()])
}
