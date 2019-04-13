// +build !amd64

package native

func DotProductComplex(input []complex64, taps []float32) complex64 {
	panic("No native function available for arch")
}

func DotProductFloat(input []float32, taps []float32) float32 {
	panic("No native function available for arch")
}

func GetNativeDotProductComplex() func(input []complex64, taps []float32) complex64 {
	return nil
}

func GetNativeDotProductFloat() func(input []float32, taps []float32) float32 {
	return nil
}

func GetNativeDotProductComplexComplex() func(input []complex64, taps []complex64) complex64 {
	return nil
}

func MultiplyConjugateInline(vecA, vecB []complex64, length int) {
	panic("No native function available for arch")
}

func GetMultiplyConjugateInline() func(vecA, vecB []complex64, length int) {
	return nil
}

func GetMultiplyConjugate() func(vecA, vecB []complex64, length int) []complex64 {
	return nil
}

func GetRotateComplex() func(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	return nil
}

func GetRotateComplexBuffer() func(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	return nil
}

func GetFirFilter() func(input []complex64, output []complex64, taps []float32) {
	return nil
}

func GetFirFilterDecimate() func(decimation uint, input []complex64, output []complex64, taps []float32) {
	return nil
}

func FirFilter(input []complex64, output []complex64, taps []float32) {
	panic("No native function available for arch")
}

func FirFilterDecimate(decimation uint, input []complex64, output []complex64, taps []float32) {
	panic("No native function available for arch")
}

func MultiplyFloatFloatVectors(A, B []float32) {
	panic("No native function available for arch")
}

func DivideFloatFloatVectors(A, B []float32) {
	panic("No native function available for arch")
}

func AddFloatFloatVectors(A, B []float32) {
	panic("No native function available for arch")
}

func SubtractFloatFloatVectors(A, B []float32) {
	panic("No native function available for arch")
}

func GetMultiplyFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetDivideFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetAddFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetSubtractFloatFloatVectors() func(A, B []float32) {
	return nil
}
