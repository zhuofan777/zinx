package main

import (
	"zinx/znet"
)

func main() {
	a := znet.NewServer("777")
	a.Serve()
}
