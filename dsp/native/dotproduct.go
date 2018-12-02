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
