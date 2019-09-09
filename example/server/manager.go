package main

import (
	"log"

	"github.com/drivebyer/zgo"
)

func main() {
	s := zgo.MakeTCPServer()
	s.Addr = "localhost"
	s.Port = "9999"
	err := s.ListenAndAccept()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	zgo.Processors[int(NetCSJoin_CODE1)][int(NetCSJoin_CODE2)] = &CSJoinProcessor
}
