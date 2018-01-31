package main

import "goio/client"

func main() {
	c := client.NewClient()
	c.Connect("172.16.6.144:10001")
}
