package digital

type DifferentialDecoder struct{}

func MakeDifferentialDecoder() *DifferentialDecoder {
	return &DifferentialDecoder{}
}

func (dd *DifferentialDecoder) Work(data []byte) []byte {
	out := make([]byte, len(data))
	dd.WorkBuffer(data, out)
	return out
}

func (dd *DifferentialDecoder) WorkBuffer(input, output []byte) int {
	lastBit := byte(0)
	mask := byte(0)

	if len(input) > len(output) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		mask = ((input[i] >> 1) & 0x7F) | (lastBit << 7)
		lastBit = input[i] & 1
		output[i] = input[i] ^ mask
	}

	return len(input)
}

func (dd *DifferentialDecoder) PredictOutputSize(inputLength int) int {
	return inputLength
}
