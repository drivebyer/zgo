package zgo

import (
	"encoding/binary"
	"log"

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
func Encode(code1 int, code2 int, stubID uint64, pb proto.Message) []byte {

	// csJoin := pb.(*message.NetCSJoin) // TODO: make more flexiable, since pb may be hold difference type.
	content, err := proto.Marshal(pb)
	if err != nil {
		log.Fatal(err)
	}
	msgLen := headerLen + len(content)

	buf := make([]byte, headerLen, msgLen)
	buf[code1Off] = byte(code1)                                                               // encode code 1
	binary.BigEndian.PutUint16(buf[code2Off:code2Off+code2Len], uint16(code2))                // encode code 2
	binary.BigEndian.PutUint32(buf[lengthOff:lengthOff+lengthLen], uint32(msgLen))            // encode massage length
	binary.BigEndian.PutUint32(buf[magicNumOff:magicNumOff+magicNumLen], uint32(magicNumber)) // encode magic number
	binary.BigEndian.PutUint64(buf[stubIDOff:stubIDOff+stubIDLen], uint64(stubID))            // encode stub ID
	binary.BigEndian.PutUint64(buf[checkSumOff:checkSumOff+checkSumLen],
		makeCheckSum(uint8(code1),
			uint16(code2),
			uint32(msgLen),
			uint32(magicNumber),
			uint64(stubID))) //encode check sum
	// fmt.Println("checksum:", makeCheckSum(uint8(message.NetCSJoin_ID_value[CODE1]),
	// 	uint16(message.NetCSJoin_ID_value[CODE1]),
	// 	uint32(msgLen),
	// 	uint32(magicNumber),
	// 	uint64(teststubID)))
	return append(buf, content...)
}
