package main

import (
	"fmt"
	"goio/inter"
)

func main() {
	var s inter.S
	s.Put()
	fmt.Println(s.Get())
}
