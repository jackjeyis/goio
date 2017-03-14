package main

import (
	"fmt"
	"runtime"
	"time"
)

func GC(sleep int) {
	var m runtime.MemStats

	for {
		time.Sleep(sleep * time.Second)
		fmt.Println("gc")
		//            runtime.GC()
		//   var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d\n", m.HeapSys, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased)

	}

}
func main() {
	GC(3)
}
