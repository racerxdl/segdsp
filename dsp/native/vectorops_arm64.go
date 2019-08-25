package native

// region Float-Float Vector Operations
func MultiplyFloatFloatVectors(A, B []float32) {
	if nativeMultiplyFloatFloatVectors == nil {
		nativeMultiplyFloatFloatVectors = GetNativeMultiplyFloatFloatVectors()
	}

	if nativeMultiplyFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeMultiplyFloatFloatVectors(A, B)
}

func DivideFloatFloatVectors(A, B []float32) {
	if nativeDivideFloatFloatVectors == nil {
		nativeDivideFloatFloatVectors = GetNativeDivideFloatFloatVectors()
	}

	if nativeDivideFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeDivideFloatFloatVectors(A, B)
}

func AddFloatFloatVectors(A, B []float32) {
	if nativeAddFloatFloatVectors == nil {
		nativeAddFloatFloatVectors = GetNativeAddFloatFloatVectors()
	}

	if nativeAddFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeAddFloatFloatVectors(A, B)
}

func SubtractFloatFloatVectors(A, B []float32) {
	if nativeSubtractFloatFloatVectors == nil {
		nativeSubtractFloatFloatVectors = GetNativeSubtractFloatFloatVectors()
	}

	if nativeSubtractFloatFloatVectors == nil {
		panic("No native function available for arch")
	}
	nativeSubtractFloatFloatVectors(A, B)
}

func GetNativeMultiplyFloatFloatVectors() func(A, B []float32) {
	// TODO

	return nil
}

func GetNativeDivideFloatFloatVectors() func(A, B []float32) {
	// TODO

	return nil
}

func GetNativeAddFloatFloatVectors() func(A, B []float32) {
	// TODO

	return nil
}

func GetNativeSubtractFloatFloatVectors() func(A, B []float32) {
	// TODO

	return nil
}

// endregion
// region Complex-Complex Vector Operations
func MultiplyComplexComplexVectors(A, B []complex64) {
	if nativeMultiplyComplexComplexVectors == nil {
		nativeMultiplyComplexComplexVectors = GetNativeMultiplyComplexComplexVectors()
	}

	if nativeMultiplyComplexComplexVectors == nil {
		panic("No native function available for arch")
	}
	nativeMultiplyComplexComplexVectors(A, B)
}

func DivideComplexComplexVectors(A, B []complex64) {
	if nativeDivideComplexComplexVectors == nil {
		nativeDivideComplexComplexVectors = GetNativeDivideComplexComplexVectors()
	}

	if nativeDivideComplexComplexVectors == nil {
		panic("No native function available for arch")
	}
	nativeDivideComplexComplexVectors(A, B)
}

func AddComplexComplexVectors(A, B []complex64) {
	if nativeAddComplexComplexVectors == nil {
		nativeAddComplexComplexVectors = GetNativeAddComplexComplexVectors()
	}

	if nativeAddComplexComplexVectors == nil {
		panic("No native function available for arch")
	}
	nativeAddComplexComplexVectors(A, B)
}

func SubtractComplexComplexVectors(A, B []complex64) {
	if nativeSubtractComplexComplexVectors == nil {
		nativeSubtractComplexComplexVectors = GetNativeSubtractComplexComplexVectors()
	}

	if nativeSubtractComplexComplexVectors == nil {
		panic("No native function available for arch")
	}
	nativeSubtractComplexComplexVectors(A, B)
}

func GetNativeMultiplyComplexComplexVectors() func(A, B []complex64) {
	// TODO

	return nil
}

func GetNativeDivideComplexComplexVectors() func(A, B []complex64) {
	// TODO

	return nil
}

func GetNativeAddComplexComplexVectors() func(A, B []complex64) {
	// TODO

	return nil
}

func GetNativeSubtractComplexComplexVectors() func(A, B []complex64) {
	// TODO

	return nil
}

// endregion
