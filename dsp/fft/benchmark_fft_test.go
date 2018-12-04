package fft

import (
	"runtime"
	"testing"
)

func BenchmarkFFT(b *testing.B) {
	b.StopTimer()

	runtime.GOMAXPROCS(runtime.NumCPU())

	N := 1 << 20
	a := make([]complex64, N)
	for i := 0; i < N; i++ {
		a[i] = complex(float32(i)/float32(N), 0)
	}

	EnsureRadix2Factors(N)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		FFT(a)
	}
}
