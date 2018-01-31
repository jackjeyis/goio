package netpoll

/*func main1() {
	poller, err := netpoll.New(nil)
	if err != nil {
		logger.Error(err)
		return
	}
	handle := func(conn net.Conn) {
		desc := netpoll.Must(netpoll.HandleRead(conn))
		poller.Start(desc, func(ev EpollEvent) {
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				poller.Stop(desc)
				return
			}
		})

	}
}
*/
