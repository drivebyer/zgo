package main

import (
	"github.com/drivebyer/zgo"
)

func main() {
	c := &zgo.Config{}
	c.Addr = "localhost"
	c.Port = "9999"
	zgo.MakeTCPServer(c)
}
