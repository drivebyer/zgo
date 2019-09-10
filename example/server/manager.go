package main

import (
	"github.com/drivebyer/zgo"
)

func main() {
	s := zgo.MakeTCPServer()
	s.Addr = ""
	s.Port = "9999"
	s.ListenAndAccept()
}

func init() {
	zgo.Processors[int(NetCSJoin_CODE1)][int(NetCSJoin_CODE2)] = &CSJoinProcessor
}
