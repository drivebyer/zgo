package zgo

import (
	"log"
	"net"
)

// RPCServer is self-explanatory.
type TCPServer struct {
	l net.Listener
}

type Config struct {
	Addr string
	Port string
}

// MakeRPCServer return a new RPCServer.
func MakeTCPServer(c *Config) *TCPServer {
	s := makeTCPServer(c)
	listen(s, c) // no block
	return s
}

// makeRPCServer return a new RPCServer.
func makeTCPServer(c *Config) *TCPServer {
	return &TCPServer{}
}

// rpcListen start listen
func listen(s *TCPServer, c *Config) {
	l, err := net.Listen("tcp", c.Addr+":"+c.Port) // TODO: make configurable
	defer l.Close()                                // close the listener when process exit.
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			return // must retrun in case
		}
		go handleConn(conn)
	}
}

// handleConn handle connection no block.
func handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		// If there is no message on the conn after 5s, we will close the conn.
		//conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		var buf [256]byte
		_, err := conn.Read(buf[:])
		if err != nil { // err != nil means that we hit the EOF of message on the conn or something else.
			log.Printf("connection over! end with err %v\n", err)
			return
		}

		decode(buf[:])
	}
}
