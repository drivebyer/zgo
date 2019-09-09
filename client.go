package zgo

import (
	"log"
	"net"
)

// TCPServer defines parameters for running an TCP server.
type TCPClient struct {
	RemotePort string
	RemoteAddr string
	activeConn map[*conn]struct{}
}

func MakeTCPClient() *TCPClient {
	return &TCPClient{}
}

func (client *TCPClient) DialAndAccept() net.Conn {
	conn, err := net.Dial("tcp", client.RemoteAddr+":"+client.RemotePort)
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		b := make([]byte, 256)
		for {
			_, err := conn.Read(b)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("Client receive %d bytes\n", n)
			Decode(b, &conn)
		}
	}()

	return conn
}
