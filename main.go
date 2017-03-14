package main

import (
	"fmt"
	"time"
	"wolfhead/timer"
)

func main() {
	timer := timer.NewTimer(1, 1000*60*60*24)
	timer_id := timer.AddTimer(3000, func() {
		fmt.Println("hello timer")
	})
	time.Sleep(60 * time.Second)
}
