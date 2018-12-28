package bitworker

// Based on libcorrect: https://github.com/quiet/libcorrect/

type BitWriter struct {
	currentByte    uint8
	currentByteLen int
	bytes          []uint8
	index          int
}

func MakeBitWriter(data []uint8) *BitWriter {
	var br = BitWriter{}
	br.Reconfigure(data)
	return &br
}

func (bw *BitWriter) Reconfigure(bytes []uint8) {
	bw.bytes = bytes
	bw.currentByte = 0
	bw.currentByteLen = 0
	bw.index = 0
}

func (bw *BitWriter) WriteN(val uint8, n int) {
	for i := 0; i < n; i++ {
		bw.Write(val)
		val >>= 1
	}
}

func (bw *BitWriter) Write(val uint8) {
	bw.currentByte |= val & 1
	bw.currentByteLen++

	if bw.currentByteLen == 8 {
		bw.bytes[bw.index] = bw.currentByte
		bw.index++
		bw.currentByte = 0
		bw.currentByteLen = 0
	} else {
		bw.currentByte <<= 1
	}
}

func (bw *BitWriter) WriteBitList(l []uint8) {
	closeLen := int(8 - bw.currentByteLen)
	if closeLen > len(l) {
		closeLen = len(l)
	}

	b := uint16(bw.currentByte)

	for i := 0; i < closeLen; i++ {
		b |= uint16(l[i])
		b <<= 1
	}

	l = l[closeLen:]

	bytes := bw.bytes
	byteIndex := bw.index

	if bw.currentByteLen+closeLen == 8 {
		b >>= 1
		bytes[byteIndex] = uint8(b & 0xFF)
		byteIndex++
	} else {
		bw.currentByte = uint8(b & 0xFF)
		bw.currentByteLen += closeLen
	}

	fullBytes := len(l) / 8

	for i := 0; i < fullBytes; i++ {
		bytes[byteIndex] = l[0]<<7 | l[1]<<6 | l[2]<<5 |
			l[3]<<4 | l[4]<<3 | l[5]<<2 |
			l[6]<<1 | l[7]
		l = l[8:]
		byteIndex++
	}

	b = 0
	for i := 0; i < len(l); i++ {
		b |= uint16(l[i])
		b <<= 1
	}

	bw.currentByte = uint8(b & 0xFF)
	bw.index = byteIndex
	bw.currentByteLen = len(l)
}

func (bw *BitWriter) WriteBitListReversed(l []uint8) {
	lr := reverse(l)
	bw.WriteBitList(lr)
}

func (bw *BitWriter) WriteFlushByte() {
	if bw.currentByteLen != 0 {
		bw.currentByte <<= uint8(8 - bw.currentByteLen)
		bw.bytes[bw.index] = bw.currentByte
		bw.index++
		bw.currentByteLen = 0
	}
}

func (bw *BitWriter) Size() int {
	return bw.index
}

func reverse(numbers []uint8) []uint8 {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
