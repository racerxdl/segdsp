package dsp

const dcFilterAlpha = 5e-5

type DCFilter struct {
	iAverage float32
	qAverage float32
}

func MakeDCFilter() *DCFilter {
	return &DCFilter{
		iAverage: 0,
		qAverage: 0,
	}
}

func (dc *DCFilter) WorkInline(data []complex64) {
	iAvg := dc.iAverage
	qAvg := dc.qAverage

	for i := 0; i < len(data); i++ {
		var s = data[i]
		iAvg = dcFilterAlpha*(real(s)-iAvg) + iAvg
		qAvg = dcFilterAlpha*(imag(s)-qAvg) + qAvg

		data[i] = complex(real(s)-iAvg, imag(s)-qAvg)
	}

	dc.iAverage = iAvg
	dc.qAverage = qAvg
}

func (dc *DCFilter) Work(data []complex64) []complex64 {
	var output = make([]complex64, len(data))
	iAvg := dc.iAverage
	qAvg := dc.qAverage

	for i := 0; i < len(data); i++ {
		var s = data[i]
		iAvg = dcFilterAlpha*(real(s)-iAvg) + iAvg
		qAvg = dcFilterAlpha*(imag(s)-qAvg) + qAvg

		output[i] = complex(real(s)-iAvg, imag(s)-qAvg)
	}

	dc.iAverage = iAvg
	dc.qAverage = qAvg

	return output
}

func (dc *DCFilter) WorkBuffer(input, output []complex64) int {
	iAvg := dc.iAverage
	qAvg := dc.qAverage

	if len(output) < len(input) {
		panic("There is not enough space in output buffer")
	}

	for i := 0; i < len(input); i++ {
		var s = input[i]
		iAvg = dcFilterAlpha*(real(s)-iAvg) + iAvg
		qAvg = dcFilterAlpha*(imag(s)-qAvg) + qAvg

		output[i] = complex(real(s)-iAvg, imag(s)-qAvg)
	}

	dc.iAverage = iAvg
	dc.qAverage = qAvg

	return len(input)
}

func (dc *DCFilter) PredictOutputSize(inputLength int) int {
	return inputLength
}
