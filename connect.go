package main

import "goio/client"

func main() {
	c := client.NewClient()
	c.Connect("127.0.0.1:18080")
}
