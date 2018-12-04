package fft

import (
	"math"
	"sync"
)

var (
	workerPoolSize = 0
	radix2Lock     sync.RWMutex
	radix2Factors  = map[int][]complex64{
		4: {complex(1, 0), complex(0, -1), complex(-1, 0), complex(0, 1)},
	}
	sqrt22 = float32(math.Sqrt2 / 2)
)

// SetWorkerPoolSize sets the number of workers during FFT computation on multicore systems.
// If n is 0 (the default), then GOMAXPROCS workers will be created.
func SetWorkerPoolSize(n int) {
	if n < 0 {
		n = 0
	}

	workerPoolSize = n
}

type fftWork struct {
	start, end int
}
