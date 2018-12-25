package dsp

import (
	"math"
)

const bwPercent = 0.80

// FloatResampler is a Polyphase Resampler based on GNU Radio Implementation
type FloatResampler struct {
	internalBuffer       []float32
	taps                 [][]float32
	diffTaps             [][]float32
	filters              []*FloatFirFilter
	diffFilters          []*FloatFirFilter
	filterSize           int
	tapsPerFilter        uint32
	decimationRate       uint32
	filterRate           float32
	lastFilter           int
	estimatedPhaseChange float32
	accumulator          float32
	rate                 float32
}

func MakeFloatResampler(filterSize int, rate float32) *FloatResampler {
	var taps []float32

	var halfBand = float64(0.5 * rate)
	var bandWidth = float64(bwPercent * halfBand)
	var transitionWidth = float64((bwPercent / 2.0) * halfBand)

	if rate >= 1 {
		halfBand = 0.5
		bandWidth = float64(bwPercent * halfBand)
		transitionWidth = float64((bwPercent / 2.0) * halfBand)
	}

	if rate < 1 {
		taps = MakeLowPassFixed(float64(filterSize), float64(filterSize), bandWidth-transitionWidth, filterSize*8)
	} else {
		taps = MakeLowPassFixed(float64(filterSize), float64(filterSize), bandWidth-transitionWidth, filterSize)
	}

	var nullTaps = make([]float32, filterSize)
	var filters = make([]*FloatFirFilter, filterSize)
	var diffFilters = make([]*FloatFirFilter, filterSize)

	for i := 0; i < filterSize; i++ {
		filters[i] = MakeFloatFirFilter(nullTaps)
		diffFilters[i] = MakeFloatFirFilter(nullTaps)
	}

	var resampler = FloatResampler{
		filterSize:     filterSize,
		lastFilter:     (len(taps) / 2) % filterSize,
		accumulator:    0,
		internalBuffer: make([]float32, 0),
		filters:        filters,
		diffFilters:    diffFilters,
		rate:           rate,
	}
	resampler.setRate(rate)
	resampler.setTaps(taps)

	var delay = float64(-filterSize) * (float64(resampler.tapsPerFilter) - 1.0) / 2.0
	var iDelay = int64(math.Round(delay))
	var acc = float64(float64(-iDelay) * float64(resampler.filterRate))
	var accInt = int64(acc)
	var accFrac = acc - float64(accInt)

	var endFilter = int64(math.Round(math.Mod(float64(resampler.lastFilter)-float64(iDelay)*float64(resampler.decimationRate)+float64(accInt), float64(filterSize))))

	resampler.estimatedPhaseChange = float32(float64(resampler.lastFilter) - (float64(endFilter) + float64(accFrac)))

	return &resampler
}

func (f *FloatResampler) setRate(rate float32) {
	f.decimationRate = uint32(math.Floor(float64(f.filterSize) / float64(rate)))
	f.filterRate = float32(float64(f.filterSize)/float64(rate) - float64(f.decimationRate))
}

func (f *FloatResampler) setTaps(taps []float32) {
	var diffTaps = generateDiffTaps(taps)

	f.taps = make([][]float32, f.filterSize)
	f.diffTaps = make([][]float32, f.filterSize)

	f.createTaps(taps, f.taps, f.filters)
	f.createTaps(diffTaps, f.diffTaps, f.diffFilters)
}

func (f *FloatResampler) createTaps(taps []float32, ourTaps [][]float32, filter []*FloatFirFilter) {
	f.tapsPerFilter = uint32(math.Ceil(float64(len(taps)) / float64(f.filterSize)))
	var tmpTaps = make([]float32, f.filterSize*int(f.tapsPerFilter))

	copy(tmpTaps, taps)

	for i := len(taps); i < len(tmpTaps); i++ {
		tmpTaps[i] = 0
	}

	for i := 0; i < int(f.filterSize); i++ {
		ourTaps[i] = make([]float32, f.tapsPerFilter)
		for j := 0; j < int(f.tapsPerFilter); j++ {
			ourTaps[i][j] = tmpTaps[i+j*f.filterSize]
		}
		filter[i].SetTaps(ourTaps[i])
	}
}

func (f *FloatResampler) filter(input []float32, length int) (int, []float32) {
	var output = make([]float32, f.PredictOutputSize(len(input)))
	var read = 0
	var wrote = 0
	var j = f.lastFilter

	for read < length {
		for j < f.filterSize {
			var o0 = f.filters[j].FilterSingle(input[read:])
			var o1 = f.diffFilters[j].FilterSingle(input[read:])

			output[wrote] = o0 + o1*f.accumulator
			wrote++

			f.accumulator += f.filterRate
			j += int(f.decimationRate) + int(math.Floor(float64(f.accumulator)))
			f.accumulator = float32(math.Mod(float64(f.accumulator), 1.0))
		}

		read += int(j / f.filterSize)
		j = j % f.filterSize
	}

	f.lastFilter = j

	output = output[:wrote]

	return read, output
}

func (f *FloatResampler) filterBuffer(input, output []float32, length int) (read, wrote int) {
	read = 0
	wrote = 0
	var j = f.lastFilter

	for read < length {
		for j < f.filterSize {
			var o0 = f.filters[j].FilterSingle(input[read:])
			var o1 = f.diffFilters[j].FilterSingle(input[read:])

			output[wrote] = o0 + o1*f.accumulator
			wrote++

			f.accumulator += f.filterRate
			j += int(f.decimationRate) + int(math.Floor(float64(f.accumulator)))
			f.accumulator = float32(math.Mod(float64(f.accumulator), 1.0))
		}

		read += int(j / f.filterSize)
		j = j % f.filterSize
	}

	f.lastFilter = j

	return read, wrote
}

func (f *FloatResampler) Work(data []float32) []float32 {
	var samples = append(f.internalBuffer, data...)

	consumed, processed := f.filter(samples, len(samples))

	if consumed < len(samples) {
		f.internalBuffer = samples[consumed:]
	} else {
		f.internalBuffer = make([]float32, 0)
	}

	return processed
}

func (f *FloatResampler) WorkBuffer(input, output []float32) int {
	if len(output) < f.PredictOutputSize(len(input)) {
		panic("There is not enough space in output buffer")
	}

	var samples = append(f.internalBuffer, input...)

	consumed, wrote := f.filterBuffer(samples, output, len(samples))

	if consumed < len(samples) {
		f.internalBuffer = samples[consumed:]
	} else {
		f.internalBuffer = make([]float32, 0)
	}

	return wrote
}

func (f *FloatResampler) PredictOutputSize(inputLength int) int {
	return int(float32(inputLength) * 2 * f.rate)
}
