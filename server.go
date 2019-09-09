package zgo

import (
	"log"
	"net"
)

// TCPServer defines parameters for running an TCP server.
type TCPServer struct {
	Addr       string
	Port       string
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
}

//
type conn struct {
	rwc net.Conn
	net.TCPConn
}

// MakeRPCServer return a new RPCServer.
func MakeTCPServer() *TCPServer {
	s := &TCPServer{}
	return s
}

// rpcListen start listen
func (s *TCPServer) ListenAndAccept() error {
	l, err := net.Listen("tcp", s.Addr+":"+s.Port) // TODO: make configurable
	defer l.Close()                                // close the listener when process exit.
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
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
		Decode(buf[:], &conn)
	}
}
