package inter

type S struct {
	i int
}

func (s *S) Get() interface{} {
	fmt.Prinln("Windows Get")
	return s.i
}

func (s *S) Put() {
	fmt.Println("Windows Put")
	s.i = 1
}
