package binarySlicer

const frequencyWindow = 1.0 / 3.0

// Float4LevelSlicer does a 4 FSK Binary Slicing assuming symbols spaced as [ -1.0, -1/3, 1/3, 1.0 ]
// It returns a byte array with symbols represented as [ 0, 1, 2, 3 ]
// It does DC Offset removal by using a simple complementary filter.
type Float4LevelSlicer struct {
	alpha   float32
	beta    float32
	average float32
}

// MakeFloat4LevelSlicer creates an instance of 4 FSK Binary Slicer
// alpha parameter represents the complementary filter strength
func MakeFloat4LevelSlicer(alpha float32) *Float4LevelSlicer {
	return &Float4LevelSlicer{
		alpha:   alpha,
		beta:    1.0 - alpha,
		average: 0,
	}
}

// Work performs the DC Offset removal and 4 level slicing
func (b4 *Float4LevelSlicer) Work(data []float32) []byte {
	var output = make([]byte, len(data))
	b4.WorkBuffer(data, output)
	return output
}

// Work performs the DC Offset removal and 4 level slicing
func (b4 *Float4LevelSlicer) WorkBuffer(input []float32, output []byte) int {
	for i := 0; i < len(input); i++ {
		var sample = input[i]

		// Recalculate DC Offset
		b4.average = b4.average*b4.beta + sample*b4.alpha

		// Remove DC Offset
		sample -= b4.average

		if sample > 0 {
			if sample >= 1.0-frequencyWindow {
				output[i] = 3
			} else {
				output[i] = 2
			}
		} else {
			if -sample >= 1.0-frequencyWindow {
				output[i] = 0
			} else {
				output[i] = 1
			}
		}
	}
	return len(input)
}

func (b4 *Float4LevelSlicer) PredictOutputSize(inputLength int) int {
	return inputLength
}
