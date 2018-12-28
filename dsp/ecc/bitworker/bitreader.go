package bitworker

var reverseTable [256]uint8

func reverseByte(b uint8) uint8 {
	return (b&0x80)>>7 | (b&0x40)>>5 | (b&0x20)>>3 |
		(b&0x10)>>1 | (b&0x08)<<1 | (b&0x04)<<3 |
		(b&0x02)<<5 | (b&0x01)<<7
}

func init() {
	for i := 0; i < 255; i++ {
		reverseTable[i] = reverseByte(uint8(i))
	}
}

type BitReader struct {
	currentByte    uint8
	currentByteLen int
	bytes          []uint8
	index          int
}

func MakeBitReader(data []uint8) *BitReader {
	var br = BitReader{}
	br.Reconfigure(data)
	return &br
}

func (br *BitReader) Reconfigure(data []uint8) {
	br.bytes = data
	br.currentByteLen = 8
	if data != nil {
		br.currentByte = data[0]
	}
	br.index = 0
}

func (br *BitReader) Read(n int) uint8 {
	read := uint8(0)
	nCopy := n

	if br.currentByteLen < n {
		read = br.currentByte & uint8((1<<uint(br.currentByteLen))-1)
		br.index++
		br.currentByte = br.bytes[br.index]
		n -= br.currentByteLen
		br.currentByteLen = 8
		read <<= uint8(n)
	}

	copyMask := (1 << uint(n)) - 1
	copyMask <<= uint(br.currentByteLen - n)
	read |= uint8(uint(int(br.currentByte)&copyMask) >> uint(br.currentByteLen-n))
	br.currentByteLen -= n

	return reverseTable[read] >> uint(8-nCopy)
}
