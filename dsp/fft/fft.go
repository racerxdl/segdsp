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

// Package fft provides forward and inverse fast Fourier transform functions.
package fft

import "github.com/racerxdl/segdsp/tools"

// FFTReal returns the forward FFT of the real-valued slice.
func FFTReal(x []float32) []complex64 {
	return FFT(tools.ToComplex64Array(x))
}

// IFFTReal returns the inverse FFT of the real-valued slice.
func IFFTReal(x []float32) []complex64 {
	return IFFT(tools.ToComplex64Array(x))
}

// IFFT returns the inverse FFT of the complex-valued slice.
func IFFT(x []complex64) []complex64 {
	lx := len(x)
	r := make([]complex64, lx)

	// Reverse inputs, which is calculated with modulo N, hence x[0] as an outlier
	r[0] = x[0]
	for i := 1; i < lx; i++ {
		r[i] = x[lx-i]
	}

	r = FFT(r)

	N := complex(float32(lx), 0)
	for n := range r {
		r[n] /= N
	}
	return r
}

// Convolve returns the convolution of x ∗ y.
func Convolve(x, y []complex64) []complex64 {
	if len(x) != len(y) {
		panic("arrays not of equal size")
	}

	fftX := FFT(x)
	fftY := FFT(y)

	r := make([]complex64, len(x))
	for i := 0; i < len(r); i++ {
		r[i] = fftX[i] * fftY[i]
	}

	return IFFT(r)
}

// FFT returns the forward FFT of the complex-valued slice.
func FFT(x []complex64) []complex64 {
	lx := len(x)

	// todo: non-hack handling length <= 1 cases
	if lx <= 1 {
		r := make([]complex64, lx)
		copy(r, x)
		return r
	}

	if tools.IsPowerOf2(lx) {
		return radix2FFT(x)
	}

	return bluesteinFFT(x)
}

// FFT2Real returns the 2-dimensional, forward FFT of the real-valued matrix.
func FFT2Real(x [][]float32) [][]complex64 {
	return FFT2(tools.ToComplex64Array2(x))
}

// FFT2 returns the 2-dimensional, forward FFT of the complex-valued matrix.
func FFT2(x [][]complex64) [][]complex64 {
	return computeFFT2(x, FFT)
}

// IFFT2Real returns the 2-dimensional, inverse FFT of the real-valued matrix.
func IFFT2Real(x [][]float32) [][]complex64 {
	return IFFT2(tools.ToComplex64Array2(x))
}

// IFFT2 returns the 2-dimensional, inverse FFT of the complex-valued matrix.
func IFFT2(x [][]complex64) [][]complex64 {
	return computeFFT2(x, IFFT)
}

func computeFFT2(x [][]complex64, fftFunc func([]complex64) []complex64) [][]complex64 {
	rows := len(x)
	if rows == 0 {
		panic("empty input array")
	}

	cols := len(x[0])
	r := make([][]complex64, rows)
	for i := 0; i < rows; i++ {
		if len(x[i]) != cols {
			panic("ragged input array")
		}
		r[i] = make([]complex64, cols)
	}

	for i := 0; i < cols; i++ {
		t := make([]complex64, rows)
		for j := 0; j < rows; j++ {
			t[j] = x[j][i]
		}

		for n, v := range fftFunc(t) {
			r[n][i] = v
		}
	}

	for n, v := range r {
		r[n] = fftFunc(v)
	}

	return r
}
