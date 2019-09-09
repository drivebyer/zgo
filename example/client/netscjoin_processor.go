package main

import (
	"fmt"

	"github.com/drivebyer/zgo"

	"github.com/golang/protobuf/proto"
)

type NetSCJoinProcessor struct {
	CODE1 int
	CODE2 int
}

var SCJoinProcessor = NetSCJoinProcessor{int(NetSCJoin_CODE1), int(NetSCJoin_CODE2)}

func (p *NetSCJoinProcessor) Handler(c *zgo.Connection, buf []byte) {
	msg := &NetSCJoin{}
	proto.Unmarshal(buf, msg)
	fmt.Println(msg.GetResp())
	//zgo.Connection.WriteAndFlush()
}
