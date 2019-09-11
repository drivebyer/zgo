package zgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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

type TLVMessage struct {
	code1    int8
	code2    int16
	length   int32
	magicNum int32
	checkSum int64
	stubID   int64
	content  []byte
}

func (msg *TLVMessage) decode(r io.Reader) {
	if err := binary.Read(r, binary.BigEndian, &msg.code1); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.code2); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.length); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.magicNum); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.checkSum); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.stubID); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, msg.content); err != nil {
		log.Fatal(err)
	}
}

func decode(buf []byte, c *Connection) {
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

	c.StubID = stubID
	c.rcvDataTime = time.Now()
	cg.add(c)

	// fmt.Println(code1, code2, reflect.TypeOf(Processors[code1][code2]), reflect.ValueOf(Processors[code1][code2]))
	var p LogicProcessor = Processors[code1][code2]

	p.Handler(c, buf[stubIDOff+stubIDLen:len(buf)])
}

func (p *TLVMessage) encode(c *Connection) {
	if err := binary.Write(c.Conn, binary.BigEndian, p.code1); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(c.Conn, binary.BigEndian, p.code2); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(c.Conn, binary.BigEndian, p.length); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(c.Conn, binary.BigEndian, p.magicNum); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(c.Conn, binary.BigEndian, p.checkSum); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(c.Conn, binary.BigEndian, p.stubID); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(c.Conn, binary.BigEndian, p.content); err != nil {
		log.Fatal(err)
	}
}

// Encode encode a TLV format which include TLV header and content
func encode(code1 int, code2 int, stubID uint64, pb proto.Message) []byte {

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

// ReadAndHandle read data from the conn, then dispatch the data to handler.
func (conn *Connection) ReadAndHandle() error {

	s := bufio.NewScanner(conn.Conn)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF {
			len := binary.BigEndian.Uint32(data[3:7])
			//binary.Read(bytes.NewReader(data[3:7]), binary.BigEndian, len)
			// fmt.Println("ReadAndHandle", len, data)
			return int(len), data[:len], nil // Note the len is the total length of message.
		}
		if atEOF {
			fmt.Println("1")
			return 0, nil, io.EOF
		}
		fmt.Println("2")
		return 0, nil, nil
	})

	var p LogicProcessor
	for s.Scan() {
		if s.Err() == io.EOF {
			os.Exit(0)
			panic(s.Err())
		}
		msg := new(TLVMessage)
		len := binary.BigEndian.Uint32(s.Bytes()[3:7])
		msg.content = make([]byte, len-headerLen)
		fmt.Println("ReadAndHandle len", len)
		msg.decode(bytes.NewReader(s.Bytes()))
		log.Println("ReadAndHandle", msg, s.Bytes())
		p = Processors[msg.code1][msg.code2]
		fmt.Println("ReadAndHandle content", msg.content)

		conn.StubID = binary.BigEndian.Uint64(s.Bytes()[19:27])
		conn.rcvDataTime = time.Now()
		cg.add(conn)

		p.Handler(conn, msg.content)
	}
	if err := s.Err(); err == nil { // err == nil, means that err == io.EOF
		fmt.Println("Error", err)
		return io.EOF
	}

	return nil
}
