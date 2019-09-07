// This file hold
package zgo

// TLV related.
const (
	magicNumber = 19700101

	// code1 code2
	lengthFieldOff = 3

	lengthFieldLen = 4

	magicNumLen = 4

	checkSumLen = 8

	stubIDLen = 8

	teststubID = 68

	// TLV header length
	headerLen = lengthFieldOff + lengthFieldLen + checkSumLen + magicNumLen + stubIDLen

	maxCode1 = 100

	maxCode2 = 2000
)

// makeCheckSum implement a simple check sum algorithm.
// TODO: adopt TCP check sum.
func makeCheckSum(code1 uint8, code2 uint16, len uint32, mn uint32, stubID uint64) uint64 {
	return uint64(code1) + uint64(code2) + uint64(len) + uint64(mn) + stubID
}
