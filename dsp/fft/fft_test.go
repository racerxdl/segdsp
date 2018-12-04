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
	"github.com/racerxdl/segdsp/tools"
	"testing"
)

// region FFT Test Helper Structs
type fftTest struct {
	in  []float32
	out []complex64
}

type reverseBitsTest struct {
	in  uint
	sz  uint
	out uint
}

type fft2Test struct {
	in  [][]float32
	out [][]complex64
}

// endregion
// region 1D FFT Test Data
var fftTests = []fftTest{
	// impulse responses
	{
		[]float32{1},
		[]complex64{complex(1, 0)},
	},
	{
		[]float32{1, 0},
		[]complex64{complex(1, 0), complex(1, 0)},
	},
	{
		[]float32{1, 0, 0, 0},
		[]complex64{complex(1, 0), complex(1, 0), complex(1, 0), complex(1, 0)},
	},
	{
		[]float32{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex64{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},

	// shifted impulse response
	{
		[]float32{0, 1},
		[]complex64{complex(1, 0), complex(-1, 0)},
	},
	{
		[]float32{0, 1, 0, 0},
		[]complex64{complex(1, 0), complex(0, -1), complex(-1, 0), complex(0, 1)},
	},
	{
		[]float32{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex64{
			complex(1, 0),
			complex(sqrt22, -sqrt22),
			complex(0, -1),
			complex(-sqrt22, -sqrt22),
			complex(-1, 0),
			complex(-sqrt22, sqrt22),
			complex(0, 1),
			complex(sqrt22, sqrt22)},
	},

	// other
	{
		[]float32{1, 2, 3, 4},
		[]complex64{
			complex(10, 0),
			complex(-2, 2),
			complex(-2, 0),
			complex(-2, -2)},
	},
	{
		[]float32{1, 3, 5, 7},
		[]complex64{
			complex(16, 0),
			complex(-4, 4),
			complex(-4, 0),
			complex(-4, -4)},
	},
	{
		[]float32{1, 2, 3, 4, 5, 6, 7, 8},
		[]complex64{
			complex(36, 0),
			complex(-4, 9.65685425),
			complex(-4, 4),
			complex(-4, 1.65685425),
			complex(-4, 0),
			complex(-4, -1.65685425),
			complex(-4, -4),
			complex(-4, -9.65685425)},
	},

	// non power of 2 lengths
	{
		[]float32{1, 0, 0, 0, 0},
		[]complex64{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},
	{
		[]float32{1, 2, 3},
		[]complex64{
			complex(6, 0),
			complex(-1.5, 0.8660254),
			complex(-1.5, -0.8660254)},
	},
	{
		[]float32{1, 1, 1},
		[]complex64{
			complex(3, 0),
			complex(0, 0),
			complex(0, 0)},
	},
}

// endregion
// region 2D FFT Test Data
var fft2Tests = []fft2Test{
	{
		[][]float32{{1, 2, 3}, {3, 4, 5}},
		[][]complex64{
			{complex(18, 0), complex(-3, 1.73205081), complex(-3, -1.73205081)},
			{complex(-6, 0), complex(0, 0), complex(0, 0)}},
	},
	{
		[][]float32{{0.1, 0.2, 0.3, 0.4, 0.5}, {1, 2, 3, 4, 5}, {3, 2, 1, 0, -1}},
		[][]complex64{
			{complex(21.5, 0), complex(-0.25, 0.34409547), complex(-0.25, 0.08122992), complex(-0.25, -0.08122992), complex(-0.25, -0.34409548)},
			{complex(-8.5, -8.66025404), complex(5.70990854, 4.6742225), complex(1.156942, 4.411357), complex(-1.65694356, 4.24889709), complex(-6.20990854, 3.98603154)},
			{complex(-8.5, 8.66025404), complex(-6.2099085, -3.9860315), complex(-1.65694356, -4.24889709), complex(1.15694356, -4.411357), complex(5.70990854, -4.6742225)}},
	},
}

// endregion
// region Reverse Bits Test Data
var reverseBitsTests = []reverseBitsTest{
	{0, 1, 0},
	{1, 2, 2},
	{1, 4, 8},
	{2, 4, 4},
	{3, 4, 12},
}

// endregion

func TestFFT(t *testing.T) {
	for _, ft := range fftTests {
		v := FFTReal(ft.in)
		if !tools.Complex64ArrayEqual(v, ft.out) {
			t.Error("FFT error\ninput:", ft.in, "\noutput:", v, "\nexpected:", ft.out)
		}

		vi := IFFT(ft.out)
		if !tools.Complex64ArrayEqual(vi, tools.ToComplex64Array(ft.in)) {
			t.Error("IFFT error\ninput:", ft.out, "\noutput:", vi, "\nexpected:", tools.ToComplex64Array(ft.in))
		}
	}
}

func TestFFT2(t *testing.T) {
	for _, ft := range fft2Tests {
		v := FFT2Real(ft.in)
		if !tools.Complex64Array2Equal(v, ft.out) {
			t.Error("FFT2 error\ninput:", ft.in, "\noutput:", v, "\nexpected:", ft.out)
		}

		vi := IFFT2(ft.out)
		if !tools.Complex64Array2Equal(vi, tools.ToComplex64Array2(ft.in)) {
			t.Error("IFFT2 error\ninput:", ft.out, "\noutput:", vi, "\nexpected:", tools.ToComplex64Array2(ft.in))
		}
	}
}

func TestReverseBits(t *testing.T) {
	for _, rt := range reverseBitsTests {
		v := reverseBits(rt.in, rt.sz)

		if v != rt.out {
			t.Error("reverse bits error\ninput:", rt.in, "\nsize:", rt.sz, "\noutput:", v, "\nexpected:", rt.out)
		}
	}
}

func TestFFTMulti(t *testing.T) {
	N := 1 << 8
	a := make([]complex64, N)
	for i := 0; i < N; i++ {
		a[i] = complex(float32(i)/float32(N), 0)
	}

	FFT(a)
}
