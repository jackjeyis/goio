package msg

type Channel interface {
	Start()
	OnRead()
	OnWrite()
	OnClose()
	DecodeMessage() error
	EncodeMessage(Message)
	Serve(Message)
	SetAttr(key, value string)
	GetAttr(string) string
}
