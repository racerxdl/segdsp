package dsp

type ComplexWorker interface {
	Work(input []complex64) []complex64
	WorkBuffer(input, output []complex64) int
	PredictOutputSize(inputLength int) int
}

type Float32Worker interface {
	Work(input []float32) []float32
	WorkBuffer(input, output []float32) int
	PredictOutputSize(inputLength int) int
}

type Complex2Float32Worker interface {
	Work(input []complex64) []float32
	WorkBuffer(input []complex64, output []float32) int
	PredictOutputSize(inputLength int) int
}
