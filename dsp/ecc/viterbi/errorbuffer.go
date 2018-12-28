package viterbi

// errorBuffer
// From https://github.com/quiet/libcorrect/blob/master/src/convolutional/error_buffer.c
type errorBuffer struct {
	index       uint32
	errors      [2][]distance
	numStates   uint32
	readErrors  []distance
	writeErrors []distance
}

func makeErrorBuffer(numStates uint32) *errorBuffer {
	eb := &errorBuffer{
		index: 0,
		errors: [2][]distance{
			make([]distance, numStates),
			make([]distance, numStates),
		},
		numStates: numStates,
	}

	eb.readErrors = eb.errors[0]
	eb.readErrors = eb.errors[1]

	return eb
}

func (eb *errorBuffer) Reset() {
	eb.errors[0] = make([]distance, eb.numStates)
	eb.errors[1] = make([]distance, eb.numStates)

	eb.index = 0

	eb.readErrors = eb.errors[0]
	eb.readErrors = eb.errors[1]
}

func (eb *errorBuffer) Swap() {
	eb.readErrors = eb.errors[eb.index]
	eb.index = (eb.index + 1) % 2
	eb.writeErrors = eb.errors[eb.index]
}
