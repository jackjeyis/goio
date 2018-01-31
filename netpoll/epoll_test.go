package netpoll

import (
	"fmt"
	"testing"
)

func TestEpollCreate(t *testing.T) {
	s, err := EpollCreate(func(e error) {
		fmt.Println(e)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}
