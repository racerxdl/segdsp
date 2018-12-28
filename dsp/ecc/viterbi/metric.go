package viterbi

import "math/bits"

// measure the hamming distance of two bit strings
// implemented as population count of x XOR y

func metricHammingDistance(x, y uint32) distance {
	return distance(bits.OnesCount32(x ^ y))
}

func metricSoftDistanceLinear(hardX uint32, softY []uint8) distance {
	dist := 0
	for i := 0; i < len(softY); i++ {
		softX := int(uint(0-int(hardX&1)) & 0xff)
		hardX >>= 1
		d := int(softY[i]) - softX
		if d < 0 {
			dist -= d
		} else {
			dist += d
		}
	}

	return distance(dist)
}

func metricSoftDistanceQuadratic(hardX uint32, softY []uint8) distance {
	dist := 0
	for i := 0; i < len(softY); i++ {
		softX := 0
		if hardX&1 > 0 {
			softX = 255
		}
		hardX >>= 1
		d := int(softY[i]) - softX
		dist += d * d
	}

	return distance(dist)
}
