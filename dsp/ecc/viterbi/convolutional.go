package viterbi

import (
	"github.com/racerxdl/segdsp/dsp/ecc/bitworker"
)

// Convolutional
// Ported from https://github.com/quiet/libcorrect/blob/master/src/convolutional/decode.c
type Convolutional struct {
	table     []uint32
	bitWriter *bitworker.BitWriter
	bitReader *bitworker.BitReader
	rate      uint32
	order     uint32
	numStates int

	hasInitDecode   bool
	distances       []distance
	pairLookup      pairLookup
	softMeasurement int

	historyBuffer *historyBuffer
	errors        *errorBuffer
}

func MakeConvolutional(rate, order uint32, poly []polynomial) *Convolutional {
	if rate < 2 {
		panic("rate must be 2 or greater")
	}

	if order > 32 {
		panic("order must be less than 32")
	}

	conv := &Convolutional{}

	conv.order = order
	conv.rate = rate
	conv.numStates = 1 << uint(order)
	conv.table = make([]uint32, 1<<order)
	conv.table = fillTable(conv.order, conv.rate, poly)

	conv.bitWriter = bitworker.MakeBitWriter(nil)
	conv.bitReader = bitworker.MakeBitReader(nil)

	conv.hasInitDecode = false

	return conv
}

func (conv *Convolutional) EncodeLength(msgLength int) int {
	msgBits := msgLength * 8
	return int(conv.rate) * (msgBits + int(conv.order) + 1)
}

func (conv *Convolutional) Encode(input []byte) []byte {
	// convolutional code convolves filter coefficients, given by
	//     the polynomial, with some history from our message.
	//     the history is stored as single subsequent bits in shiftregister
	sr := shiftRegister(0)

	// shiftMask is the sr bit mask that removes bits
	//      that extend beyond order
	// e.g. if order is 7, then remove the 8th bit and beyond
	shiftMask := shiftRegister((1 << conv.order) - 1)

	encodedLenBits := conv.EncodeLength(len(input))
	encodedLength := encodedLenBits / 8
	if encodedLenBits%8 > 0 {
		encodedLength = encodedLenBits/8 + 1
	}

	outBuf := make([]byte, encodedLength)

	conv.bitWriter.Reconfigure(outBuf)
	conv.bitReader.Reconfigure(input)

	for i := 0; i < 8*len(input); i++ {
		// shiftregister has oldest bits on left, newest on right
		sr <<= 1
		sr |= shiftRegister(conv.bitReader.Read(1))
		sr &= shiftMask
		// shift most significant bit from byte and move down one bit at a time

		// we do direct lookup of our convolutional output here
		// all of the bits from this convolution are stored in this row
		out := conv.table[sr]
		conv.bitWriter.WriteN(uint8(out&0xFF), int(conv.rate))
	}

	// now flush the shiftregister
	// this is simply running the loop as above but without any new inputs
	// or rather, the new input string is all 0s
	for i := uint32(0); i < conv.order+1; i++ {
		sr <<= 1
		sr &= shiftMask
		out := conv.table[sr]
		conv.bitWriter.WriteN(uint8(out&0xFF), int(conv.rate))
	}

	// 0-fill any remaining bits on our final byte
	conv.bitWriter.WriteFlushByte()

	return outBuf
}

func (conv *Convolutional) Decode(input []byte) []byte {
	if uint32(len(input))%conv.rate > 0 {
		panic("Input length should be a multiple of rate")
	}

	numEncodedBytes := len(input) / 8
	if len(input)%8 > 0 {
		numEncodedBytes = 1 + len(input)/8
	}

	conv.bitReader.Reconfigure(input)

	return conv.decodeWork(input, numEncodedBytes, false)
}

func (conv *Convolutional) DecodeSoft(input []byte) []byte {
	if uint32(len(input))%conv.rate > 0 {
		panic("Input length should be a multiple of rate")
	}

	numEncodedBytes := len(input) / 8
	if len(input)%8 > 0 {
		numEncodedBytes = 1 + len(input)/8
	}

	return conv.decodeWork(input, numEncodedBytes, true)
}

func (conv *Convolutional) decodeInit(minTraceback, tracebackLength, renormalizeInterval uint32) {
	conv.hasInitDecode = true
	conv.distances = make([]distance, 1<<conv.rate)
	conv.pairLookup = pairLookupCreate(conv.rate, conv.order, conv.table)
	conv.softMeasurement = CorrectSoftLinear
	// we limit history to go back as far as 5 * the order of our polynomial
	conv.historyBuffer = makeHistoryBuffer(minTraceback, tracebackLength, renormalizeInterval, uint32(conv.numStates/2), 1<<(conv.order-1))
	conv.errors = makeErrorBuffer(uint32(conv.numStates))
}

func (conv *Convolutional) decodeWork(input []byte, encodedBytes int, soft bool) []byte {
	if !conv.hasInitDecode {
		maxErrorPerInput := conv.rate * softMax
		renormalizeInterval := distanceMax / uint32(maxErrorPerInput)
		conv.decodeInit(5*uint32(conv.order), 15*uint32(conv.order), renormalizeInterval)
	}

	sets := uint32(len(input)) / conv.rate

	outBuf := make([]byte, encodedBytes)

	conv.bitWriter.Reconfigure(outBuf)

	conv.errors.Reset()
	conv.historyBuffer.Reset()

	conv.warmup(sets, input, soft)
	conv.inner(sets, input, soft)
	conv.tail(sets, input, soft)

	conv.historyBuffer.Flush(conv.bitWriter)

	return outBuf
}

func (conv *Convolutional) warmup(sets uint32, input []byte, soft bool) {
	// first phase: load shiftregister up from 0 (order goes from 1 to conv->order)
	// we are building up error metrics for the first order bits
	for i := uint32(0); i < conv.order-1 && i < sets; i++ {
		out := uint32(0)
		if !soft {
			out = uint32(conv.bitReader.Read(int(conv.rate)))
		}

		readErrors := conv.errors.readErrors
		writeErrors := conv.errors.writeErrors

		// walk all of the state we have so far
		for j := 0; j < 1<<(i+1); j++ {
			last := j >> 1
			var dist distance
			if soft {
				v := input[i*conv.rate:]
				if conv.softMeasurement == CorrectSoftLinear {
					dist = metricSoftDistanceLinear(conv.table[j], v)
				} else {
					dist = metricSoftDistanceQuadratic(conv.table[j], v)
				}
			} else {
				dist = metricHammingDistance(conv.table[j], out)
			}
			writeErrors[j] = dist + readErrors[last]
		}
		conv.errors.Swap()
	}
}

func (conv *Convolutional) inner(sets uint32, input []byte, soft bool) {
	highBit := shiftRegister(1 << (conv.order - 1))

	for i := conv.order - 1; i < (sets - conv.order + 1); i++ {
		distances := conv.distances
		// lasterrors are the aggregate bit errors for the states of shiftregister for the previous
		// time slice

		if soft {
			v := input[i*conv.rate : (i+1)*conv.rate]
			if conv.softMeasurement == CorrectSoftLinear {
				for j := uint32(0); j < 1<<conv.rate; j++ {
					distances[j] = metricSoftDistanceLinear(j, v)
				}
			} else {
				for j := uint32(0); j < 1<<conv.rate; j++ {
					distances[j] = metricSoftDistanceQuadratic(j, v)
				}
			}
		} else {
			out := uint32(conv.bitReader.Read(int(conv.rate)))
			for j := uint32(0); j < 1<<conv.rate; j++ {
				distances[j] = metricHammingDistance(j, out)
			}
		}

		pairLookupFillDistances(&conv.pairLookup, distances)

		// a mask to get the high order bit from the shift register
		numIter := highBit << 1
		readErrors := conv.errors.readErrors
		// aggregate bit errors for this time slice
		writeErrors := conv.errors.writeErrors

		history := conv.historyBuffer.GetSlice()
		// walk through all states, ignoring oldest bit
		// we will track a best register state (path) and the number of bit errors at that path at
		// this time slice
		// this loop considers two paths per iteration (high order bit set, clear)
		// so, it only runs numstates/2 iterations
		// we'll update the history for every state and find the path with the least aggregated bit
		// errors

		// now run the main loop
		// we calculate 2 sets of 2 register states here (4 states per iter)
		// this creates 2 sets which share a predecessor, and 2 sets which share a successor
		//
		// the first set definition is the two states that are the same except for the least order
		// bit
		// these two share a predecessor because their high n - 1 bits are the same (differ only by
		// newest bit)
		//
		// the second set definition is the two states that are the same except for the high order
		// bit
		// these two share a successor because the oldest high order bit will be shifted out, and
		// the other bits will be present in the successor
		//

		highBase := int(highBit >> 1)
		low := 0
		high := highBit
		base := 0
		for high < numIter {
			// shifted-right ancestors
			// low and low_plus_one share low_past_error
			//   note that they are the same when shifted right by 1
			// same goes for high and high_plus_one

			offset := 0
			baseOffset := 0
			for baseOffset < 4 {
				lowKey := conv.pairLookup.keys[base+baseOffset]
				highKey := conv.pairLookup.keys[highBase+baseOffset]
				lowConcatDist := conv.pairLookup.distances[lowKey]
				highConcatDist := conv.pairLookup.distances[highKey]

				lowPastError := readErrors[base+baseOffset]
				highPastError := readErrors[highBase+base+baseOffset]

				lowError := (lowConcatDist & 0xFFFF) + lowPastError
				highError := (highConcatDist & 0xFFFF) + highPastError

				successor := low + offset

				error_ := highError
				historyMask := uint8(1)

				if lowError <= highError {
					error_ = lowError
					historyMask = 0
				}

				writeErrors[successor] = error_
				history[successor] = historyMask

				lowPlusOneError := (uint32(lowConcatDist) >> 16) + uint32(lowPastError)
				highPlusOneError := (uint32(highConcatDist) >> 16) + uint32(highPastError)

				plusOneSuccesor := lowPlusOneError

				plusOneError := highPlusOneError
				plusOneHistoryMask := uint8(1)

				if lowPlusOneError <= highPlusOneError {
					plusOneError = lowPlusOneError
					plusOneHistoryMask = 0
				}

				writeErrors[plusOneSuccesor] = distance(plusOneError)
				history[plusOneSuccesor] = plusOneHistoryMask

				offset += 2
				baseOffset++
			}

			low += 8
			high += 8
			base += 4
		}

		conv.historyBuffer.Process(writeErrors, conv.bitWriter)
		conv.errors.Swap()
	}
}

func (conv *Convolutional) tail(sets uint32, input []byte, soft bool) {
	// flush state registers
	// now we only shift in 0s, skipping 1-successors

	highBit := 1 << (conv.order - 1)

	for i := sets - conv.order + 1; i < sets; i++ {
		// lasterrors are the aggregate bit errors for the states of shiftregister for the previous
		// time slice
		writeErrors := conv.errors.writeErrors
		readErrors := conv.errors.readErrors
		history := conv.historyBuffer.GetSlice()
		distances := conv.distances

		// calculate the distance from all output states to our sliced bits
		if soft {
			v := input[i*conv.rate : (i+1)*conv.rate]
			if conv.softMeasurement == CorrectSoftLinear {
				for j := uint32(0); j < 1<<conv.rate; j++ {
					distances[j] = metricSoftDistanceLinear(j, v)
				}
			} else {
				for j := uint32(0); j < 1<<conv.rate; j++ {
					distances[j] = metricSoftDistanceQuadratic(j, v)
				}
			}
		} else {
			out := uint32(conv.bitReader.Read(int(conv.rate)))
			for j := uint32(0); j < 1<<conv.rate; j++ {
				distances[j] = metricHammingDistance(j, out)
			}
		}

		table := conv.table

		// a mask to get the high order bit from the shift register
		numIter := highBit << 1
		skip := 1 << (conv.order - (sets - i))
		baseSkip := skip >> 1

		highBase := highBit >> 1
		low := 0
		high := highBit
		base := 0

		for high < numIter {
			lowOutput := table[low]
			highOutput := table[high]

			lowDist := distances[lowOutput]
			highDist := distances[highOutput]

			lowPastError := readErrors[base]
			highPastError := readErrors[highBase+base]

			lowError := lowDist + lowPastError
			highError := highDist + highPastError

			successor := low

			error_ := highError
			historyMask := uint8(1)

			if lowError < highError {
				error_ = lowError
				historyMask = 0
			}

			writeErrors[successor] = error_
			history[successor] = historyMask

			low += skip
			high += skip
			base += baseSkip
		}

		conv.historyBuffer.ProcessSkip(writeErrors, conv.bitWriter, uint32(skip))
		conv.errors.Swap()
	}
}
