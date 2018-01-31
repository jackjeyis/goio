package netpoll

import (
	"fmt"
	"goio/netpoll"
	"testing"
)

func TestKqueueCreate(t *testing.T) {
	c := &netpoll.KqueueConfig{}
	k, err := netpoll.KqueueCreate(c)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(k)
}
