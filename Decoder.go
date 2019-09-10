package zgo

import (
	"encoding/binary"
	"log"
	"net"
	"time"
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

func decode(buf []byte, conn *net.Conn) {
	code1 := uint8(buf[code1Off])
	code2 := binary.BigEndian.Uint16(buf[code2Off : code2Off+code2Len])

	length := binary.BigEndian.Uint32(buf[lengthOff : lengthOff+lengthLen])

	// 0x12345678
	// lower addr -------> higher addr
	// 0x12  |  0x34  |  0x56  |  0x78
	magicNum := binary.BigEndian.Uint32(buf[magicNumOff : magicNumOff+magicNumLen])
	if magicNum != magicNumber {
		log.Fatal("MagicNum Fial")
	}

	stubID := binary.BigEndian.Uint64(buf[stubIDOff : stubIDOff+stubIDLen])

	checkSum := binary.BigEndian.Uint64(buf[checkSumOff : checkSumOff+checkSumLen])

	// fmt.Println("Before return decode", buf[0], buf[1:3], buf[3:7], buf[7:11], buf[11:19], buf[19:27])
	// fmt.Println("Before return decode", buf[code1Off],
	// 	buf[code2Off:code2Off+code2Len],
	// 	buf[lengthOff:lengthOff+lengthLen],
	// 	buf[magicNumOff:magicNumOff+magicNumLen],
	// 	buf[checkSumOff:checkSumOff+checkSumLen],
	// 	buf[stubIDOff:stubIDOff+stubIDLen])
	// fmt.Println("Before return decode", code1, code2, length, magicNum, checkSum, stubID)
	// fmt.Println("server check sum:", makeCheckSum(code1, code2, length, magicNumber, stubID))
	if makeCheckSum(code1, code2, length, magicNumber, stubID) != checkSum {
		log.Fatal("CheckSum Fail")
	}

	c := &Connection{}
	c.StubID = stubID
	c.rcvDataTime = time.Now()
	c.Conn = *conn
	cg.add(c)

	// fmt.Println(code1, code2, reflect.TypeOf(Processors[code1][code2]), reflect.ValueOf(Processors[code1][code2]))
	var p LogicProcessor = Processors[code1][code2]

	p.Handler(c, buf[stubIDOff+stubIDLen:len(buf)])
}
