package zgo

import (
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

type connPool struct {
	connections map[*Connection]struct{}
}

var cg connPool

func (c connPool) add(con *Connection) {
	c.connections[con] = struct{}{}
}

func (c connPool) remove(conn *Connection) {
	delete(c.connections, conn)
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
	// log.Println(c.Conn, code1, code2, stubID, msg)
	if _, err := c.Conn.Write(b); err != nil {
		log.Fatal(err)
	}
}
