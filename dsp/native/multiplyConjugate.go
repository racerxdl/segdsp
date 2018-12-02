// +build !amd64

package native

func MultiplyConjugate(vecA, vecB []complex64) {
	panic("No native function available for arch")
}

func GetMultiplyConjugate() func(vecA, vecB []complex64) {
	return nil
}

func MultiplyConjugate(vecA, vecB []complex64, length int) []complex64 {
	panic("No native function available for arch")
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
