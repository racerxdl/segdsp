package viterbi

import (
	"github.com/racerxdl/segdsp/dsp/ecc/bitworker"
	"math"
)

// historyBuffer is a ring buffer of path histories
// generates output bits after accumulating sufficient history
// Ported from https://github.com/quiet/libcorrect/blob/master/src/convolutional/history_buffer.c
type historyBuffer struct {
	// history entries must be at least this old to be decoded
	minTracebackLength uint32

	// we'll decode entries in bursts. this tells us the length of the burst
	tracebackGroupLength uint32

	// we will store a total of cap entries. equal to min_traceback_length +
	// traceback_group_length
	capacity uint32

	// how many states in the shift register? this is one of the dimensions of
	// history table
	numStates uint32

	// what's the high order bit of the shift register?
	highBit shiftRegister

	// history is a compact history representation for every shift register
	// state,
	//    one bit per time slice
	history [][]uint8

	// which slice are we writing next?
	index uint32

	// how many valid entries are there?
	length uint32

	// temporary store of fetched bits
	fetched []uint8

	// how often should we renormalize?
	renormalizeInterval uint32
	renormalizeCounter  uint32
}

func makeHistoryBuffer(minTracebackLength, tracebackGroupLength, renormalizeInterval, numStates uint32, highBit shiftRegister) *historyBuffer {
	hb := &historyBuffer{
		minTracebackLength:   minTracebackLength,
		tracebackGroupLength: tracebackGroupLength,
		capacity:             minTracebackLength + tracebackGroupLength,
		numStates:            numStates,
		highBit:              highBit,
	}

	hb.history = make([][]uint8, hb.capacity)
	for i := range hb.history {
		hb.history[i] = make([]uint8, hb.numStates)
	}

	hb.fetched = make([]uint8, hb.capacity)

	hb.index = 0
	hb.length = 0
	hb.renormalizeCounter = 0
	hb.renormalizeInterval = renormalizeInterval

	return hb
}

func (hb *historyBuffer) Step() {

}

func (hb *historyBuffer) GetSlice() []uint8 {
	return hb.history[hb.index]
}

func (hb *historyBuffer) Reset() {
	hb.index = 0
	hb.length = 0
}

func (hb *historyBuffer) Search(distances []distance, every uint32) shiftRegister {
	var bestPath shiftRegister
	leastError := distance(math.MaxUint16)
	// search for a state with the least error
	for i := uint32(0); i < hb.numStates; i += every {
		if distances[i] < leastError {
			leastError = distances[i]
			bestPath = shiftRegister(i)
		}
	}

	return bestPath
}

func (hb *historyBuffer) Renormalize(distances []distance, minRegister shiftRegister) {
	minDistance := distances[minRegister]
	for i := uint32(0); i < hb.numStates; i++ {
		distances[i] -= minDistance
	}
}

func (hb *historyBuffer) TraceBack(bestPath shiftRegister, minTracebackLength uint32, bitWriter *bitworker.BitWriter) {
	fetchedIndex := 0
	highBit := hb.highBit
	index := hb.index
	capacity := hb.capacity

	for j := uint32(0); j < minTracebackLength; j++ {
		if index == 0 {
			index = capacity - 1
		} else {
			index--
		}

		// we're walking backwards from what the work we did before
		// so, we'll shift high order bits in
		// the path will cross multiple different shift register states, and we determine
		//   which state by going backwards one time slice at a time

		history := hb.history[index][bestPath]
		pathBit := highBit
		if history > 0 {
			pathBit = 0
		}

		bestPath |= pathBit
		bestPath >>= 1
	}

	prefetchIndex := index

	if prefetchIndex == 0 {
		prefetchIndex = capacity - 1
	} else {
		prefetchIndex--
	}

	length := hb.length

	for j := minTracebackLength; j < length; j++ {
		index = prefetchIndex
		if prefetchIndex == 0 {
			prefetchIndex = capacity - 1
		} else {
			prefetchIndex--
		}

		//prefetch(buf->history[prefetch_index]);

		// we're walking backwards from what the work we did before
		// so, we'll shift high order bits in
		// the path will cross multiple different shift register states, and we determine
		//   which state by going backwards one time slice at a time
		history := hb.history[index][bestPath]
		pathBit := highBit
		if history > 0 {
			pathBit = 0
		}

		bestPath |= pathBit
		bestPath >>= 1
		hb.fetched[fetchedIndex] = 1
		if pathBit > 0 {
			hb.fetched[fetchedIndex] = 0
		}

		fetchedIndex++
	}

	bitWriter.WriteBitListReversed(hb.fetched)

	hb.length -= uint32(fetchedIndex)
}

func (hb *historyBuffer) ProcessSkip(distances []distance, output *bitworker.BitWriter, skip uint32) {
	hb.index++
	if hb.index == hb.capacity {
		hb.index = 0
	}

	hb.renormalizeCounter++
	hb.length++

	// there are four ways these branches can resolve
	// a) we are neither renormalizing nor doing a traceback
	// b) we are renormalizing but not doing a traceback
	// c) we are renormalizing and doing a traceback
	// d) we are not renormalizing but we are doing a traceback
	// in case c, we want to save the effort of finding the bestpath
	//    since that's expensive
	// so we have to check for that case after we renormalize
	if hb.renormalizeCounter == hb.renormalizeInterval {
		hb.renormalizeCounter = 0
		bestPath := hb.Search(distances, skip)
		hb.Renormalize(distances, bestPath)

		if hb.length == hb.capacity {
			// reuse the bestpath found for renormalizing
			hb.TraceBack(bestPath, hb.minTracebackLength, output)
		}
	} else if hb.length == hb.capacity {
		bestPath := hb.Search(distances, skip)
		hb.TraceBack(bestPath, hb.minTracebackLength, output)
	}
}

func (hb *historyBuffer) Process(distances []distance, output *bitworker.BitWriter) {
	hb.ProcessSkip(distances, output, 1)
}

func (hb *historyBuffer) Flush(output *bitworker.BitWriter) {
	hb.TraceBack(0, 0, output)
}
