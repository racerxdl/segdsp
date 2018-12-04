/*
 * Copyright (c) 2011 Matt Jibson <matt.jibson@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package fft

import (
	"math"
	"runtime"
	"sync"
)

// EnsureRadix2Factors ensures that all radix 2 factors are computed for inputs
// of length input_len. This is used to precompute needed factors for known
// sizes. Generally should only be used for benchmarks.
func EnsureRadix2Factors(length int) {
	getRadix2Factors(length)
}

func getRadix2Factors(length int) []complex64 {
	radix2Lock.RLock()

	if hasRadix2Factors(length) {
		defer radix2Lock.RUnlock()
		return radix2Factors[length]
	}

	radix2Lock.RUnlock()
	radix2Lock.Lock()
	defer radix2Lock.Unlock()

	if !hasRadix2Factors(length) {
		for i, p := 8, 4; i <= length; i, p = i<<1, i {
			if radix2Factors[i] == nil {
				radix2Factors[i] = make([]complex64, i)

				for n, j := 0, 0; n < i; n, j = n+2, j+1 {
					radix2Factors[i][n] = radix2Factors[p][j]
				}

				for n := 1; n < i; n += 2 {
					sin, cos := math.Sincos(float64(-2 * math.Pi / float32(i) * float32(n)))
					radix2Factors[i][n] = complex(float32(cos), float32(sin))
				}
			}
		}
	}

	return radix2Factors[length]
}

func hasRadix2Factors(idx int) bool {
	return radix2Factors[idx] != nil
}

// radix2FFT returns the FFT calculated using the radix-2 DIT Cooley-Tukey algorithm.
func radix2FFT(x []complex64) []complex64 {
	lx := len(x)
	factors := getRadix2Factors(lx)

	t := make([]complex64, lx) // temp
	r := reorderData(x)

	var blocks, stage, s2 int

	jobs := make(chan *fftWork, lx)
	wg := sync.WaitGroup{}

	numWorkers := workerPoolSize
	if (numWorkers) == 0 {
		numWorkers = runtime.GOMAXPROCS(0)
	}

	idxDiff := lx / numWorkers
	if idxDiff < 2 {
		idxDiff = 2
	}

	worker := func() {
		for work := range jobs {
			for nb := work.start; nb < work.end; nb += stage {
				if stage != 2 {
					for j := 0; j < s2; j++ {
						idx := j + nb
						idx2 := idx + s2
						ridx := r[idx]
						w_n := r[idx2] * factors[blocks*j]
						t[idx] = ridx + w_n
						t[idx2] = ridx - w_n
					}
				} else {
					n1 := nb + 1
					rn := r[nb]
					rn1 := r[n1]
					t[nb] = rn + rn1
					t[n1] = rn - rn1
				}
			}
			wg.Done()
		}
	}

	for i := 0; i < numWorkers; i++ {
		go worker()
	}
	defer close(jobs)

	for stage = 2; stage <= lx; stage <<= 1 {
		blocks = lx / stage
		s2 = stage / 2

		for start, end := 0, stage; ; {
			if end-start >= idxDiff || end == lx {
				wg.Add(1)
				jobs <- &fftWork{start, end}

				if end == lx {
					break
				}

				start = end
			}

			end += stage
		}
		wg.Wait()
		r, t = t, r
	}

	return r
}

// reorderData returns a copy of x reordered for the DFT.
func reorderData(x []complex64) []complex64 {
	lx := uint(len(x))
	r := make([]complex64, lx)
	s := log2(lx)

	var n uint
	for ; n < lx; n++ {
		r[reverseBits(n, s)] = x[n]
	}

	return r
}

// log2 returns the log base 2 of v
// from: http://graphics.stanford.edu/~seander/bithacks.html#IntegerLogObvious
func log2(v uint) uint {
	var r uint

	for v >>= 1; v != 0; v >>= 1 {
		r++
	}

	return r
}

// reverseBits returns the first s bits of v in reverse order
// from: http://graphics.stanford.edu/~seander/bithacks.html#BitReverseObvious
func reverseBits(v, s uint) uint {
	var r uint

	// Since we aren't reversing all the bits in v (just the first s bits),
	// we only need the first bit of v instead of a full copy.
	r = v & 1
	s--

	for v >>= 1; v != 0; v >>= 1 {
		r <<= 1
		r |= v & 1
		s--
	}

	return r << s
}
