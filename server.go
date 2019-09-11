package zgo

import (
	"log"
	"net"
	"time"
)

// TCPServer defines parameters for running an TCP server.
type TCPServer struct {
	Addr      string
	Port      string
	listeners map[*net.Listener]struct{}
	//activeConn map[*conn]struct{}

	gp     *gpool
	gpSize int32
}

// MakeRPCServer return a new RPCServer.
func MakeTCPServer() *TCPServer {
	s := &TCPServer{}
	s.gp = MakeGPool(s.gpSize) // Default pool size.
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
		conn := s.accept(l)

		// Every time we accept a conneciton, we try to reuse a goroutine in
		// gpool. Only When The gpool if empty, we create a goroutine.
		if !s.gp.Get(conn) {
			go func(conn Connection) {
				for {
					err := conn.ReadAndHandle()
					// var buf [32]byte
					// _, err := conn.Conn.Read(buf[:])
					log.Println("ListenAndAccept", conn)
					if err != nil { // err != nil means that we hit the EOF of message on the conn or something else.
						log.Printf("connection %v over! end with err %v\n", conn, err)
						conn.Close()
						i, ok := s.gp.Put()
						if !ok {
							return
						}
						conn = i.(Connection)
						continue
					}
					//decode(buf[:], &conn)

					time.Sleep(1 * time.Second)

				}
			}(conn)
		}
	}
}

func (s *TCPServer) accept(l net.Listener) Connection {
	c, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	return Connection{Conn: c}
}

//
// Although you can use SetDeadline(or else related) on every connection to handler
// the half-close situation in TCP. But there’s a balance to consider: you shouldn’t do it too early,
// in case the client is just slow on generating data. How to choosing a timeout is a problem.
// Maybe Ping the client is a better choice. This is what TCP keepalive doing.
//
// TCP keepalive. https://tools.ietf.org/html/rfc1122#page-101
func EnableTCPKeepAlive(c *net.TCPConn, time time.Duration) {
	if err := c.SetKeepAlive(true); err != nil {
		log.Fatal(err)
	}
	if err := c.SetKeepAlivePeriod(time); err != nil {
		log.Fatal(err)
	}
}
