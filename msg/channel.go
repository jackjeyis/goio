package msg

type Channel interface {
	Start()
	OnRead()
	OnWrite()
	OnClose()
	DecodeMessage() error
	EncodeMessage(Message)
	Service() Service
}
