package viterbi

import "math/bits"

type pairLookup struct {
	keys        []distancePairKey
	outputs     []outputPair
	outputMask  outputPair
	outputWidth uint32
	distances   []distance
}

func fillTable(order, rate uint32, poly []polynomial) []uint32 {
	var length = 1 << order
	var table = make([]uint32, length)

	for i := 0; i < length; i++ {
		out := 0
		mask := 1
		for j := uint32(0); j < rate; j++ {
			v := bits.OnesCount16(uint16(i)&uint16(poly[j])) % 2
			if v > 0 {
				out |= mask
			}
			mask <<= 1
		}
		table[i] = uint32(out)
	}

	return table
}

func pairLookupCreate(rate, order uint32, table []uint32) pairLookup {
	pairs := pairLookup{
		keys:    make([]distancePairKey, 1<<(order-1)),
		outputs: make([]outputPair, 1<<(rate*2)),
	}

	invOutputs := make([]uint, 1<<(rate*2))
	outputCounter := uint(1)

	l := 1 << (order - 1)

	for i := 0; i < l; i++ {
		out := table[i*2+1]
		out <<= rate
		out |= table[i*2]

		if invOutputs[out] == 0 {
			invOutputs[out] = outputCounter
			pairs.outputs[outputCounter] = outputPair(out)
			outputCounter++
		}

		pairs.keys[i] = distancePairKey(invOutputs[out])
	}

	pairs.outputs = pairs.outputs[:outputCounter]
	pairs.outputMask = (1 << rate) - 1
	pairs.outputWidth = rate
	pairs.distances = make([]distance, len(pairs.outputs))

	return pairs
}

func pairLookupFillDistances(pairs *pairLookup, distances []distance) {
	for i := 1; i < len(pairs.outputs); i++ {
		concatOut := pairs.outputs[i]
		i0 := concatOut & pairs.outputMask
		concatOut >>= pairs.outputWidth
		i1 := concatOut

		pairs.distances[i] = (distances[i0] << 16) | distances[i1]
	}
}
