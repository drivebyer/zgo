package main

import (
	"fmt"

	"github.com/drivebyer/zgo"

	"github.com/golang/protobuf/proto"
)

type NetCSJoinProcessor struct {
	CODE1 int
	CODE2 int
}

var CSJoinProcessor = NetCSJoinProcessor{int(NetCSJoin_CODE1), int(NetCSJoin_CODE2)}

func (p *NetCSJoinProcessor) Handler(c *zgo.Connection, buf []byte) {
	msg := &NetCSJoin{}
	proto.Unmarshal(buf, msg)
	fmt.Println("Get client message:", msg.GetReq(), c.StubID)

	responseMSG := NetSCJoin{}
	responseMSG.Resp = "pong"
	c.WriteAndFlush(int(NetSCJoin_CODE1), int(NetSCJoin_CODE2), c.StubID, &responseMSG)
}
