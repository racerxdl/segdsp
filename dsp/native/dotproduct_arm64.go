package native

// import "github.com/racerxdl/segdsp/dsp/native/arm64"

var nativeDotProductFloat func(input []float32, taps []float32) float32
var nativeDotProductComplex func(input []complex64, taps []float32) complex64

func DotProductComplex(input []complex64, taps []float32) complex64 {
	if nativeDotProductComplex == nil {
		nativeDotProductComplex = GetNativeDotProductComplex()
	}

	if nativeDotProductComplex == nil {
		panic("No native function available for arch")
	}
	return nativeDotProductComplex(input, taps)
}

func DotProductFloat(input []float32, taps []float32) float32 {
	if nativeDotProductFloat == nil {
		nativeDotProductFloat = GetNativeDotProductFloat()
	}

	if nativeDotProductFloat == nil {
		panic("No native function available for arch")
	}
	return nativeDotProductFloat(input, taps)
}

func GetNativeDotProductComplex() func(input []complex64, taps []float32) complex64 {
	// Neon is always available at AArch64
	// Disabled for now because we don't have AArch64 support on c2goasm
	// return amd64.DotProductNativeComplexNeon
	return nil
}

func GetNativeDotProductFloat() func(input []float32, taps []float32) float32 {
	// Neon is always available at AArch64
	// Disabled for now because we don't have AArch64 support on c2goasm
	// return arm64.DotProductNativeFloatNeon
	return nil
}
