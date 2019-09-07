package zgo

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/drivebyer/zgo/example/message"
	"github.com/golang/protobuf/proto"
)

//  TLV Header, total length: 27 bytes.
//  0 1 3 7 11 19
//	///////////////////////////////////////////////////////////////////////////////////////////
//	+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
//  +Byte Length+ 1	  		+ 2			+ 4	     	+ 4			+ 8		 	+ 8			+
//	+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
//	+Field Name + CODE1 	+ CODE2 	+ Length 	+ MagicNum  + CheckSum  + StubID	+
//	+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
//
//	CODE1 		direction: receive or send,
//	CODE2 		service type
//	Length 		length
//	MagicNum	magic number
//	CheckSum	check sum
//	StubID		stud id
//	///////////////////////////////////////////////////////////////////////////////////////////

func decode(buf []byte) {
	code1 := uint8(buf[0])
	code2 := binary.BigEndian.Uint16(buf[1:3])

	length := binary.BigEndian.Uint32(buf[3:7])

	// 0x12345678
	// lower addr -------> higher addr
	// 0x12  |  0x34  |  0x56  |  0x78
	magicNum := binary.BigEndian.Uint32(buf[7:11])
	if magicNum != magicNumber {
		log.Fatal("MagicNum Fial")
	}

	stubID := binary.BigEndian.Uint64(buf[19:27])

	checkSum := binary.BigEndian.Uint64(buf[11:19])
	if makeCheckSum(code1, code2, length, magicNum, stubID) != checkSum {
		log.Fatal("CheckSum Fail")
	}

	fmt.Println("Before return decode", buf[0], buf[1:3], buf[3:7], buf[7:11], buf[11:19], buf[19:27])
	fmt.Println("Before return decode", code1, code2, length, magicNum, checkSum, stubID)
	msg := &message.CSJoin{}
	proto.Unmarshal(buf[27:len(buf)], msg)
	fmt.Println(msg.CODE1, msg.CODE2, msg.GetReq())
}
