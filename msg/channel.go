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
	SetAttr(key, value string)
	GetAttr(string) string
}
