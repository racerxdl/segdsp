// viterbi is a Viterbi / Convolutional Encode / Decode flow in Go
// It is a ported version from quiet/libcorrect: https://github.com/quiet/libcorrect/
package viterbi

import "math"

const (
	CorrectSoftLinear = iota
	CorrectSoftQuadratic
)

const (
	softMax     = math.MaxUint8
	distanceMax = math.MaxUint16
)

type distance uint16
type shiftRegister uint16
type distancePairKey uint32
type outputPair uint32
type polynomial uint16
