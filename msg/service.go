package msg

type Service interface {
	Serve(Message)
}
