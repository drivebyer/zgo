package zgo

import "time"

// TLV related.
const (
	CODE1 = "CODE1"
	CODE2 = "CODE2"

	magicNumber = 19700101

	code1Off = 0
	code1Len = 1

	// maxCode1 = 100
	// maxCode2 = 2000

	maxCode1 = 10
	maxCode2 = 20

	code2Off = code1Off + code1Len // 1
	code2Len = 2

	lengthOff = code2Off + code2Len // 3
	lengthLen = 4

	magicNumOff = lengthOff + lengthLen // 7
	magicNumLen = 4

	checkSumOff = magicNumOff + magicNumLen // 11
	checkSumLen = 8

	stubIDOff = checkSumOff + checkSumLen // 19
	stubIDLen = 8

	// TLV header length
	headerLen = code1Len + code2Len + lengthLen + magicNumLen + stubIDLen + stubIDLen

	// How many times we try when we dont's receive data from client
	tryCount = 10
	// How long we wait before we send Ping after we dont's receive data from client
	idleTime        = 15000
	idleTimeDration = time.Duration(idleTime * time.Millisecond)
	//
	retryInterval         = 1000
	retryIntervalDutation = time.Duration(retryInterval * time.Millisecond)
	// keepalive tick interval
	tickTime         = 1000
	tickTimeDuration = time.Duration(tickTime * time.Millisecond)
)

// logicProcessors hold all logic processor which define by user.
var Processors [maxCode1][maxCode2]LogicProcessor

// makeCheckSum implement a simple check sum algorithm.
// TODO: adopt TCP check sum.
func makeCheckSum(code1 uint8, code2 uint16, len uint32, mn uint32, stubID uint64) uint64 {
	return uint64(code1) + uint64(code2) + uint64(len) + uint64(mn) + stubID
}
