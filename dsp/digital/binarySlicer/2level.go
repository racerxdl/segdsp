package binarySlicer

// Float2LevelSlicer is a Binary Slicer for 2 FSK.
// It basically returns a byte with value 0 for any float sample that is below or equal zero
// and byte with value 1 for any float sample that is over 0
type Float2LevelSlicer struct{}

// MakeFloat2LevelSlicer creates an instance of 2 FSK Binary Slicer
func MakeFloat2LevelSlicer() *Float2LevelSlicer {
	return &Float2LevelSlicer{}
}

// Work processes a FM Demodulated Float sample array to extract the binary symbols
func (b2 *Float2LevelSlicer) Work(data []float32) []byte {
	var output = make([]byte, len(data))
	b2.WorkBuffer(data, output)
	return output
}

// Work performs the DC Offset removal and 4 level slicing
func (b2 *Float2LevelSlicer) WorkBuffer(input []float32, output []byte) int {
	for i := 0; i < len(input); i++ {
		if input[i] > 0 {
			output[i] = 1
		} else {
			output[i] = 0
		}
	}
	return len(input)
}

func (b2 *Float2LevelSlicer) PredictOutputSize(inputLength int) int {
	return inputLength
}
