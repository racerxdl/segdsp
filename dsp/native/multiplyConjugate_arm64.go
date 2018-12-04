package native

// import "github.com/racerxdl/segdsp/dsp/native/arm64"

var nativeMultiplyConjugate func(vecA, vecB []complex64, length int) []complex64
var nativeMultiplyConjugateInline func(vecA, vecB []complex64, length int)

func MultiplyConjugate(vecA, vecB []complex64, length int) []complex64 {
	if nativeMultiplyConjugate == nil {
		nativeMultiplyConjugate = GetMultiplyConjugate()
	}

	if nativeMultiplyConjugate == nil {
		panic("No native function available for arch")
	}
	return nativeMultiplyConjugate(vecA, vecB, length)
}

func MultiplyConjugateInline(vecA, vecB []complex64, length int) {
	if nativeMultiplyConjugateInline == nil {
		nativeMultiplyConjugateInline = GetMultiplyConjugateInline()
	}

	if nativeMultiplyConjugateInline == nil {
		panic("No native function available for arch")
	}
	nativeMultiplyConjugateInline(vecA, vecB, length)
}

func GetMultiplyConjugateInline() func(vecA, vecB []complex64, length int) {
	// Neon is always available at AArch64
	// Disabled for now because we don't have AArch64 support on c2goasm
	// return arm64.MultiplyConjugateInlineNeon
	return nil
}

func GetMultiplyConjugate() func(vecA, vecB []complex64, length int) []complex64 {
	// Neon is always available at AArch64
	// Disabled for now because we don't have AArch64 support on c2goasm
	// return arm64.MultiplyConjugateNeon
	return nil
}
