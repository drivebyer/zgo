package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/drivebyer/zgo"
)

func main() {

	playerID, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	c := zgo.MakeTCPClient()
	c.RemoteAddr = "localhost"
	c.RemotePort = "9999"
	conn := c.DialAndAccept()
	zconn := zgo.Connection{}
	zconn.Conn = conn
	for {
		msg := NetCSJoin{}
		msg.Req = "ping"

		zconn.WriteAndFlush(int(NetCSJoin_CODE1), int(NetCSJoin_CODE2), playerID, &msg)
		time.Sleep(1 * time.Second)
	}
	// ch := make(chan os.Signal, 1)
	// select {
	// case s := <-ch:
	// 	fmt.Println("exit with signal", s)
	// }
}

func init() {
	zgo.Processors[int(NetSCJoin_CODE1)][int(NetSCJoin_CODE2)] = &SCJoinProcessor
}
