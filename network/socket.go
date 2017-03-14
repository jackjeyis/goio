package network

import (
	"goio/protocol"
	"net"
)

type Socket struct {
	conn     *net.TCPConn
	protocol protocol.Protocol
}
