package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	var (
		c     *zk.Conn
		paths []string
		err   error
	)
	c, _, err = zk.Connect([]string{"172.16.6.77"}, 500*time.Millisecond)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("connect success")
	paths, _, err = c.Children("/barrage")
	/*	paths, _, err = c.Children("/barrage/server")

		fmt.Println(paths)
		if err != nil {
			fmt.Println(err)
		}
		sort.Strings(paths)
		fmt.Println(paths)
		data, _, err := c.Get("/barrage/server/" + paths[0])
		if err != nil {
			fmt.Println(err)
		}
	*/
	//	is_exist, _, err := c.Exists("/barrage/master")
	//if err != nil {
	//fmt.Println("Exists", err)
	//}
	/*	if !is_exist {
		_, err := c.Create("/barrage/master", data, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("Create", err)
		}
	} else {*/
	//	_, stat, err := c.Get("/barrage/master")
	//_, err = c.Set("/barrage/master", data, stat.Version)
	data, _, err := c.Get("/barrage/" + paths[0])
	if err != nil {
		fmt.Println("Get", err)
	}
	fmt.Println(string(data))
	//}

	//	time.Sleep(1 * time.Hour)
}
