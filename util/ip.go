package util

import (
	"goio/logger"
	"net"
)

func InternalIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		logger.Error("net.InterfaceAddrs error (%v)", err)
		panic(err)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}
