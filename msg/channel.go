package msg

type Channel interface {
	Start()
	OnRead()
	OnWrite()
	OnClose()
	Close()
	DecodeMessage() error
	EncodeMessage(Message)
	Serve(Message)
	SetAttr(string, interface{})
	GetAttr(string) interface{}
	GetIOService() Service
	SetDeadline(int)
}
