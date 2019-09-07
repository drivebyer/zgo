package zgo

import (
	"encoding/binary"
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
//	CODE1 		direction, receive or send
//	CODE2 		service type
//	Length 		message length
//	MagicNum	magic number
//	CheckSum	check sum
//	StubID		stud id
//	///////////////////////////////////////////////////////////////////////////////////////////

type tlvMessage struct {
	code1 int
	code2 int
	msg   proto.Message // TODO: should related with proto
}

// Encode encode a TLV format which include TLV header and content
func Encode(pb proto.Message) []byte {

	csJoin := pb.(*message.CSJoin) // TODO: make more flexiable, since pb may be hold difference type.
	content, err := proto.Marshal(pb)
	if err != nil {
		log.Fatal(err)
	}
	msgLen := headerLen + len(content)
	buf := make([]byte, headerLen, msgLen)
	buf[0] = byte(csJoin.GetCODE1())                                // encode code 1
	binary.BigEndian.PutUint16(buf[1:3], uint16(csJoin.GetCODE2())) // encode code 2
	binary.BigEndian.PutUint32(buf[3:7], uint32(msgLen))            // encode massage length
	binary.BigEndian.PutUint32(buf[7:11], uint32(magicNumber))      // encode magic number
	binary.BigEndian.PutUint64(buf[19:27], uint64(teststubID))      // encode stub ID
	binary.BigEndian.PutUint64(buf[11:19],
		makeCheckSum(uint8(csJoin.GetCODE1()),
			uint16(csJoin.GetCODE2()),
			uint32(msgLen),
			uint32(magicNumber),
			uint64(teststubID))) //encode check sum

	return append(buf, content...)
}
