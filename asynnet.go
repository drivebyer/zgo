package main

import (
	"log"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "localhost:6868")
	if err != nil {
		log.Fatal(err)
	}

	// using loop to keep accepting new connections util the server died.
	for {
		// when no conn, block here. After Accept(), we may wanna do something with the conn.
		// If we wanna build no block server,
		// the should no block between the Accept() and the end of for statement.
		// Never block in the Accecp() loop.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

// Now, we can do something with the conn
// Add timeout.
func handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		// If there is no message on the conn after 5s, we will close the conn.
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var buf [256]byte
		n, err := conn.Read(buf[:])
		if err != nil { // err != nil means that we hit the EOF of message on the conn or something else.
			log.Printf("connection over! end with err %v\n", err)
			return
		}
		os.Stderr.Write(buf[:n])
	}
}
