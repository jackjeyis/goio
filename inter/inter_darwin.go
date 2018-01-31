package inter

import "fmt"

type S struct {
	i int
}

func (s *S) Get() interface{} {
	fmt.Println("darwin Get")
	return s.i
}

func (s *S) Put() {
	fmt.Println("darwin Put")
	s.i = 1
}
