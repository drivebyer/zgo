package zgo

import (
	"log"
	"net"
)

// TCPServer defines parameters for running an TCP server.
type TCPClient struct {
	RemotePort string
	RemoteAddr string
	//activeConn map[*conn]struct{}
}

func MakeTCPClient() *TCPClient {
	return &TCPClient{}
}

func (client *TCPClient) DialAndAccept() *Connection {
	c, err := net.Dial("tcp", client.RemoteAddr+":"+client.RemotePort)
	conn := &Connection{Conn: c}

	if err != nil {
		log.Fatal("DialAndAccept", err)
	}
	go func() {
		defer conn.Close()
		b := make([]byte, 256)
		for {
			_, err := conn.Conn.Read(b)
			if err != nil {
				log.Fatal("DialAndAccept", err)
			}
			// fmt.Printf("Client receive %d bytes\n", n)
			decode(b, conn)
		}
	}()

	return conn
}
