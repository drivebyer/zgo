package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/drivebyer/zgo"
	"github.com/drivebyer/zgo/example/message"
)

func main() {

	c := &zgo.Config{}
	c.Addr = "localhost"
	c.Port = "9999"
	conn, err := net.Dial("tcp", c.Addr+":"+c.Port)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg := &message.CSJoin{} // CSJoin is a type which satisfy Message.
		msg.CODE1 = 111
		msg.CODE2 = 222
		r := bufio.NewReader(os.Stdin)
		s, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		msg.Req = s
		fmt.Println("Before encode msg:", msg)
		b := zgo.Encode(msg)
		fmt.Println("After encode byte buffer:", b)
		conn.Write(b)
	}
}
