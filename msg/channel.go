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
	SetAttr(string, string)
	GetAttr(string) string
	GetIOService() Service
	SetDeadline(int)
}
