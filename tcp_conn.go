package zgo

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type connPool struct {
	mu          *sync.Mutex
	connections map[*Connection]struct{}
}

var cg connPool

func (c connPool) add(con *Connection) {
	//c.mu.Lock()
	c.connections[con] = struct{}{}
	//c.mu.Unlock()
}

func (c connPool) remove(conn *Connection) {
	//c.mu.Lock()
	delete(c.connections, conn)
	//c.mu.Unlock()
}

type Connection struct {
	Conn        net.Conn
	StubID      uint64    // once connect, bind the stubID with the connection.
	tryCount    int       // keep alive try count
	rcvDataTime time.Time // last time receive data
}

func init() {
	cg.connections = make(map[*Connection]struct{}, 256)
}

func (c *Connection) WriteAndFlush(code1 int, code2 int, stubID uint64, msg proto.Message) {
	b := encode(code1, code2, stubID, msg)
	log.Println("WriteAndFlush", code1, code2, stubID, msg)
	if _, err := c.Conn.Write(b); err != nil {
		log.Fatal(err)
	}
}

func (c *Connection) Close() {
	if err := c.Conn.Close(); err != nil {
		log.Fatal(err)
	}
	//cg.remove(c)
}

func EnableAPPKeepAlive() {
	go func() {

		for {
			cg.mu.Lock()
			for k, _ := range cg.connections {
				if time.Since(k.rcvDataTime) > time.Duration(idleTimeDration) {

				}
			}
			cg.mu.Lock()

			time.Sleep(tickTimeDuration)
		}

	}()
}
