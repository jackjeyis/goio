package main

import (
	"fmt"
	"goio/util"
)

func main() {
	var ctrie *util.Ctrie = util.New(nil)
	ctrie.Insert([]byte("key"), 1)

	a, _ := ctrie.Remove([]byte("key"))
	fmt.Println(a)
}
