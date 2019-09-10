package zgo

import (
	"log"
	"net"
)

// TCPServer defines parameters for running an TCP server.
type TCPServer struct {
	Addr      string
	Port      string
	listeners map[*net.Listener]struct{}
	//activeConn map[*conn]struct{}

	gp *gpool
}

//
// type conn struct {
// 	rwc net.Conn
// 	net.TCPConn
// }

// MakeRPCServer return a new RPCServer.
func MakeTCPServer() *TCPServer {
	s := &TCPServer{}
	s.gp = MakeGPool(10) // Default pool size.
	return s
}

// rpcListen start listen
func (s *TCPServer) ListenAndAccept() {
	l, err := net.Listen("tcp", s.Addr+":"+s.Port) // TODO: make configurable
	defer l.Close()                                // close the listener when process exit.
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Every time we accept a conneciton, we try to reuse a goroutine in
		// gpool. Only When The gpool if empty, we create a goroutine.
		if !s.gp.Get(conn) {
			go func(conn net.Conn) {
				for {
					// If there is no message on the conn after 5s, we will close the conn.
					//conn.SetReadDeadline(time.Now().Add(30 * time.Second))
					var buf [256]byte
					_, err := conn.Read(buf[:])
					log.Println("Conn ", conn)
					if err != nil { // err != nil means that we hit the EOF of message on the conn or something else.
						log.Printf("connection %v over! end with err %v\n", conn, err)
						conn.Close()

						ok, i := s.gp.Put()
						if !ok {
							return // If the pool is full, we exit the goroutine.
						}
						conn = i.(net.Conn)
						continue
					}
					decode(buf[:], &conn)
				}
			}(conn)
		}
	}
}
