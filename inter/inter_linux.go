package inter

type S struct {
	i int
}

func (s *S) Get() interface{} {
	fmt.Prinln("linux Get")
	return s.i
}

func (s *S) Put() {
	fmt.Println("linux Put")
	s.i = 1
}
