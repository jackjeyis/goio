package inter

type Store interface {
	Put()
	Get() interface{}
}
