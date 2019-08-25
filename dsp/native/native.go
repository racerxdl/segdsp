// +build !amd64,!arm64

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

func GetNativeMultiplyConjugateInline() func(vecA, vecB []complex64, length int) {
	return nil
}

func GetNativeMultiplyConjugate() func(vecA, vecB []complex64, length int) []complex64 {
	return nil
}

func GetNativeRotateComplex() func(input []complex64, phase *complex64, phaseIncrement complex64, length int) []complex64 {
	return nil
}

func GetNativeRotateComplexBuffer() func(input, output []complex64, phase *complex64, phaseIncrement complex64, length int) int {
	return nil
}

func GetNativeFirFilter() func(input []complex64, output []complex64, taps []float32) {
	return nil
}

func GetNativeFirFilterDecimate() func(decimation uint, input []complex64, output []complex64, taps []float32) {
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

func GetNativeMultiplyFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetNativeDivideFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetNativeAddFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetNativeSubtractFloatFloatVectors() func(A, B []float32) {
	return nil
}

func GetNativeMultiplyComplexComplexVectors() func(A, B []complex64) {
	return nil
}

func GetNativeDivideComplexComplexVectors() func(A, B []complex64) {
	return nil
}

func GetNativeAddComplexComplexVectors() func(A, B []complex64) {
	return nil
}

func GetNativeSubtractComplexComplexVectors() func(A, B []complex64) {
	return nil
}
