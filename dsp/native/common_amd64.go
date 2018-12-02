package native

import "github.com/racerxdl/segdsp/dsp/native/amd64"

func GetSIMDMode() string {
	if amd64.AVX {
		return "AVX"
	}
	if amd64.SSE2 {
		return "SSE2"
	}

	return "NONE"
}
