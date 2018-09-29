package dsp

import "math"

// Default hamming window

func hammingWindow(nTaps int) []float64 {
	var taps = make([]float64, nTaps)

	var M = float64(nTaps - 1)

	for i := 0; i < nTaps; i++ {
		taps[i] = 0.5 - 0.5*math.Cos((2*math.Pi*float64(i))/M)
	}

	return taps
}

func computeNTaps(sampleRate, transitionWidth float64) int {
	var maxAttenuation = 53.0
	var nTaps = int(maxAttenuation * sampleRate / (22.0 * transitionWidth))

	nTaps |= 1

	return nTaps
}

func computeNTapsAtt(sampleRate, transitionWidth, maxAttenuation float64) int {
	var nTaps = int(maxAttenuation * sampleRate / (22.0 * transitionWidth))

	nTaps |= 1

	return nTaps
}

func MakeRRC(gain, sampleRate, symbolRate, alpha float64, nTaps int) []float32 {
	nTaps |= 1
	var taps = make([]float32, nTaps)
	var spb = sampleRate / symbolRate
	var scale = float64(0)
	var x1, x2, x3, num, den, xindx float64

	for i := 0; i < nTaps; i++ {
		xindx = float64(i) - float64(nTaps)/2.0
		x1 = math.Pi * xindx / spb
		x2 = 4 * alpha * xindx / spb
		x3 = x2*x2 - 1

		if math.Abs(x3) > 0.000001 {
			if i != nTaps/2 {
				num = math.Cos((1+alpha)*x1) + math.Sin((1-alpha)*x1)/(4*alpha*xindx/spb)
			} else {
				num = math.Cos((1+alpha)*x1) + (1-alpha)*math.Pi/(4*alpha)
			}

			den = x3 * math.Pi
		} else {
			if alpha == 1 {
				taps[i] = -1
				continue
			}

			x3 = (1 - alpha) * x1
			x2 = (1 + alpha) * x1

			num = math.Sin(x2)*(1+alpha)*math.Pi - math.Cos(x3)*((1-alpha)*math.Pi*spb)/(4*alpha*xindx) + math.Sin(x3)*spb*spb/(4*alpha*xindx*xindx)
			den = -32 * math.Pi * alpha * alpha * xindx / spb
		}

		taps[i] = float32(4 * alpha * num / den)
		scale += float64(taps[i])
	}

	for i := 0; i < nTaps; i++ {
		taps[i] = float32(float64(taps[i]) * gain / scale)
		if taps[i] > 1 {
			taps[i] = 1
		}
	}

	return taps
}

func MakeLowPass(gain, sampleRate, cutFrequency, transitionWidth float64) []float32 {
	var nTaps = computeNTaps(sampleRate, transitionWidth)
	var taps = make([]float32, nTaps)
	var w = hammingWindow(nTaps)

	var M = (nTaps - 1) / 2
	var fwT0 = 2 * math.Pi * cutFrequency / sampleRate

	for i := -M; i <= M; i++ {
		if i == 0 {
			taps[i+M] = float32(fwT0 / math.Pi * w[i+M])
		} else {
			taps[i+M] = float32(math.Sin(float64(i)*fwT0) / (float64(i) * math.Pi) * w[i+M])
		}
	}

	var fmax = float64(taps[M])
	for i := 1; i <= M; i++ {
		fmax += 2 * float64(taps[i+M])
	}

	gain /= fmax

	for i := 0; i < nTaps; i++ {
		taps[i] = float32(float64(taps[i]) * gain)
	}

	return taps
}

func MakeLowPass2(gain, sampleRate, cutFrequency, transitionWidth, attenuation float64) []float32 {
	var nTaps = computeNTapsAtt(sampleRate, transitionWidth, attenuation)
	var taps = make([]float32, nTaps)
	var w = hammingWindow(nTaps)

	var M = (nTaps - 1) / 2
	var fwT0 = 2 * math.Pi * cutFrequency / sampleRate

	for i := -M; i <= M; i++ {
		if i == 0 {
			taps[i+M] = float32(fwT0 / math.Pi * w[i+M])
		} else {
			taps[i+M] = float32(math.Sin(float64(i)*fwT0) / (float64(i) * math.Pi) * w[i+M])
		}
	}

	var fmax = float64(taps[M])
	for i := 1; i <= M; i++ {
		fmax += 2 * float64(taps[i+M])
	}

	gain /= fmax

	for i := 0; i < nTaps; i++ {
		taps[i] = float32(float64(taps[i]) * gain)
	}

	return taps
}

func generateDiffTaps(taps []float32) []float32 {
	var dF0 = float64(-1)
	var dF1 = float64(1)

	var diffTaps = make([]float32, len(taps))

	for i := 0; i < len(taps)-1; i++ {
		diffTaps[i] = float32(dF0*float64(taps[i]) + dF1*float64(taps[i+1]))
	}

	diffTaps[len(taps)-1] = 0

	return diffTaps
}

func MakeLowPassFixed(gain, sampleRate, cutFrequency float64, length int) []float32 {
	length |= 1

	var taps = make([]float32, length)
	var frequency = cutFrequency / sampleRate
	var center = int(math.Floor(float64(length) / 2))
	var sum = 0.0

	for i := 0; i < length; i++ {
		var val = float64(0.0)
		if i == center {
			val = 2 * math.Pi * float64(frequency)
		} else {
			var angle = 2 * math.Pi * (float64(i) + 1) / (float64(length) + 1)
			val = math.Sin(2*math.Pi*frequency*float64(i-center)) / float64(i-center)
			val *= 0.42 - 0.5*math.Cos(angle) + 0.08*math.Cos(2*angle)
		}
		sum += val
		taps[i] = float32(val)
	}

	for i := 0; i < length; i++ {
		taps[i] /= float32(sum)
		taps[i] *= float32(gain)
	}

	return taps
}
