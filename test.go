package main

import (
	"fmt"
	"wolfhead/util"
)

func main() {
	pool := util.Pool{
		New: func() interface{} {
			return "hello"
		},
	}
	s, ok := pool.Get().(string)
	fmt.Println(s)
}
