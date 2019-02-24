package tools

import (
	"log"
	"math"
)

// All Functions here are from Standard Go Library but ported to float32

const (
	mask  = 0x7FF
	shift = 32 - 8 - 1
	bias  = 255
	//signMask = 1 << 31
	//fracMask = 1<<shift - 1

	epsilon = 1e-6
)

/*
	Floating-point arctangent.
*/

// The original C code, the long comment, and the constants below were
// from http://netlib.sandia.gov/cephes/cmath/atan.c, available from
// http://www.netlib.org/cephes/cmath.tgz.
// The go code is a version of the original C.
//
// atan.c
// Inverse circular tangent (arctangent)
//
// SYNOPSIS:
// double x, y, atan();
// y = atan( x );
//
// DESCRIPTION:
// Returns radian angle between -pi/2 and +pi/2 whose tangent is x.
//
// Range reduction is from three intervals into the interval from zero to 0.66.
// The approximant uses a rational function of degree 4/5 of the form
// x + x**3 P(x)/Q(x).
//
// ACCURACY:
//                      Relative error:
// arithmetic   domain    # trials  peak     rms
//    DEC       -10, 10   50000     2.4e-17  8.3e-18
//    IEEE      -10, 10   10^6      1.8e-16  5.0e-17
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov

// sin coefficients
var _sin = [...]float32{
	1.58962301576546568060E-10, // 0x3de5d8fd1fd19ccd
	-2.50507477628578072866E-8, // 0xbe5ae5e5a9291f5d
	2.75573136213857245213E-6,  // 0x3ec71de3567d48a1
	-1.98412698295895385996E-4, // 0xbf2a01a019bfdf03
	8.33333333332211858878E-3,  // 0x3f8111111110f7d0
	-1.66666666666666307295E-1, // 0xbfc5555555555548
}

// cos coefficients
var _cos = [...]float32{
	-1.13585365213876817300E-11, // 0xbda8fa49a0861a9b
	2.08757008419747316778E-9,   // 0x3e21ee9d7b4e3f05
	-2.75573141792967388112E-7,  // 0xbe927e4f7eac4bc6
	2.48015872888517045348E-5,   // 0x3efa01a019c844f5
	-1.38888888888730564116E-3,  // 0xbf56c16c16c14f91
	4.16666666666665929218E-2,   // 0x3fa555555555554b
}

// xatan evaluates a series valid in the range [0, 0.66].
func xatan(x float32) float32 {
	const (
		P0 = -8.750608600031904122785e-01
		P1 = -1.615753718733365076637e+01
		P2 = -7.500855792314704667340e+01
		P3 = -1.228866684490136173410e+02
		P4 = -6.485021904942025371773e+01
		Q0 = +2.485846490142306297962e+01
		Q1 = +1.650270098316988542046e+02
		Q2 = +4.328810604912902668951e+02
		Q3 = +4.853903996359136964868e+02
		Q4 = +1.945506571482613964425e+02
	)
	z := x * x
	z = z * ((((P0*z+P1)*z+P2)*z+P3)*z + P4) / (((((z+Q0)*z+Q1)*z+Q2)*z+Q3)*z + Q4)
	z = x*z + x
	return z
}

// satan reduces its argument (known to be positive)
// to the range [0, 0.66] and calls xatan.
func satan(x float32) float32 {
	const (
		Morebits = 6.123233995736765886130e-17 // pi/2 = PIO2 + Morebits
		Tan3pio8 = 2.41421356237309504880      // tan(3*pi/8)
	)
	if x <= 0.66 {
		return xatan(x)
	}
	if x > Tan3pio8 {
		return math.Pi/2 - xatan(1/x) + Morebits
	}
	return math.Pi/4 + xatan((x-1)/(x+1)) + 0.5*Morebits
}

// Atan returns the arctangent, in radians, of x.
//
// Special cases are:
//      Atan(±0) = ±0
//      Atan(±Inf) = ±Pi/2
func Atan(x float32) float32 {
	if x == 0 {
		return x
	}
	if x > 0 {
		return satan(x)
	}
	return -satan(-x)
}

// Hypot returns Sqrt(p*p + q*q), taking care to avoid
// unnecessary overflow and underflow.
//
// Special cases are:
//	Hypot(±Inf, q) = +Inf
//	Hypot(p, ±Inf) = +Inf
//	Hypot(NaN, q) = NaN
//	Hypot(p, NaN) = NaN
func Hypot(p, q float32) float32 {
	// special cases
	switch {
	case IsInf(p, 0) || IsInf(q, 0):
		return float32(math.Inf(1))
	case IsNaN(p) || IsNaN(q):
		return float32(math.NaN())
	}
	if p < 0 {
		p = -p
	}
	if q < 0 {
		q = -q
	}
	if p < q {
		p, q = q, p
	}
	if p == 0 {
		return 0
	}
	q = q / p
	return p * float32(math.Sqrt(float64(1+q*q)))
}

func ComplexAbs(x complex64) float32 { return Hypot(real(x), imag(x)) }

func ComplexAbsSquared(x complex64) float32 {
	return real(x)*real(x) + imag(x)*imag(x)
}

func Conj(x complex64) complex64 { return complex(real(x), -imag(x)) }

func IsNaN(f float32) bool {
	// IEEE 754 says that only NaNs satisfy f != f.
	// To avoid the floating-point hardware, could use:
	//	x := Float64bits(f);
	//	return uint32(x>>shift)&mask == mask && x != uvinf && x != uvneginf
	return f != f
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf(f float32, sign int) bool {
	// Test for infinity by comparing against maximum float.
	// To avoid the floating-point hardware, could use:
	//	x := Float64bits(f);
	//	return sign >= 0 && x == uvinf || sign <= 0 && x == uvneginf;
	return sign >= 0 && f > math.MaxFloat32 || sign <= 0 && f < -math.MaxFloat32
}

func Signbit(x float32) bool {
	return math.Float32bits(x)&(1<<31) != 0
}

func Copysign(x, y float32) float32 {
	const sign = 1 << 31
	return math.Float32frombits(math.Float32bits(x)&^sign | math.Float32bits(y)&sign)
}

func Atan2(y, x float32) float32 {
	// special cases
	switch {
	case IsNaN(y) || IsNaN(x):
		return float32(math.NaN())
	case y == 0:
		if x >= 0 && !Signbit(x) {
			return Copysign(0, y)
		}
		return Copysign(math.Pi, y)
	case x == 0:
		return Copysign(math.Pi/2, y)
	case IsInf(x, 0):
		if IsInf(x, 1) {
			switch {
			case IsInf(y, 0):
				return Copysign(math.Pi/4, y)
			default:
				return Copysign(0, y)
			}
		}
		switch {
		case IsInf(y, 0):
			return Copysign(3*math.Pi/4, y)
		default:
			return Copysign(math.Pi, y)
		}
	case IsInf(y, 0):
		return Copysign(math.Pi/2, y)
	}

	// Call atan and determine the quadrant.
	q := Atan(y / x)
	if x < 0 {
		if q <= 0 {
			return q + math.Pi
		}
		return q - math.Pi
	}
	return q
}

func Abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func Clip(v, max float32) float32 {
	return 0.5 * (Abs(v+max) - Abs(v-max))
}

func Modf(f float32) (int float32, frac float32) {
	if f < 1 {
		switch {
		case f < 0:
			int, frac = Modf(-f)
			return -int, -frac
		case f == 0:
			return f, f // Return -0, -0 when f == -0
		}
		return 0, f
	}

	x := math.Float32bits(f)
	e := uint(x>>shift)&mask - bias

	// Keep the top 8+e bits, the integer part; clear the rest.
	if e < 32-8 {
		x &^= 1<<(32-8-e) - 1
	}
	int = math.Float32frombits(x)
	frac = f - int
	return
}

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//	Floor(±0) = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN) = NaN
func Floor(x float32) float32 {
	if x == 0 || IsNaN(x) || IsInf(x, 0) {
		return x
	}
	if x < 0 {
		d, fract := Modf(-x)
		if fract != 0.0 {
			d = d + 1
		}
		return -d
	}
	d, _ := Modf(x)
	return d
}

func AlmostFloatEqual(a, b float32) bool {
	return Abs(a-b) <= epsilon || Abs(1-a/b) <= epsilon
}

func ComplexEqual(a, b complex64) bool {
	// Safe Compare two complexes within
	// This is needed here because multiplying using SIMD might generate slightly different value
	return AlmostFloatEqual(real(a), real(b)) && AlmostFloatEqual(imag(a), imag(b))
}

func Complex64ArrayEqual(a, b []complex64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if !ComplexEqual(a[i], b[i]) {
			log.Println("Error at", i)
			log.Printf("Expected %v got %v\n", a[i], b[i])
			return false
		}
	}

	return true
}
func Complex64Array2Equal(a, b [][]complex64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if !Complex64ArrayEqual(a[i], b[i]) {
			log.Println("A Error at", i)
			return false
		}
	}

	return true
}

func ComplexPhase(c complex64) float32 {
	return Atan2(imag(c), real(c))
}

func ComplexToPolar(c complex64) (r, θ float32) {
	return ComplexAbs(c), ComplexPhase(c)
}

func ComplexNormalize(c complex64) complex64 {
	abs := ComplexAbs(c)
	return complex(real(c)/abs, imag(c)/abs)
}

func Sin(x float32) float32 {
	const (
		PI4A = 7.85398125648498535156E-1                             // 0x3fe921fb40000000, Pi/4 split into three parts
		PI4B = 3.77489470793079817668E-8                             // 0x3e64442d00000000,
		PI4C = 2.69515142907905952645E-15                            // 0x3ce8469898cc5170,
		M4PI = 1.273239544735162542821171882678754627704620361328125 // 4/pi
	)
	// special cases
	switch {
	case x == 0 || IsNaN(x):
		return x // return ±0 || NaN()
	case IsInf(x, 0):
		return float32(math.NaN())
	}

	// make argument positive but save the sign
	sign := false
	if x < 0 {
		x = -x
		sign = true
	}

	j := int32(x * M4PI) // integer part of x/(Pi/4), as integer for tests on the phase angle
	y := float32(j)      // integer part of x/(Pi/4), as float

	// map zeros to origin
	if j&1 == 1 {
		j++
		y++
	}
	j &= 7 // octant modulo 2Pi radians (360 degrees)
	// reflect in x axis
	if j > 3 {
		sign = !sign
		j -= 4
	}

	z := ((x - y*PI4A) - y*PI4B) - y*PI4C // Extended precision modular arithmetic
	zz := z * z
	if j == 1 || j == 2 {
		y = 1.0 - 0.5*zz + zz*zz*((((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5])
	} else {
		y = z + z*zz*((((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5])
	}
	if sign {
		y = -y
	}
	return y
}

func Cos(x float32) float32 {
	const (
		PI4A = 7.85398125648498535156E-1                             // 0x3fe921fb40000000, Pi/4 split into three parts
		PI4B = 3.77489470793079817668E-8                             // 0x3e64442d00000000,
		PI4C = 2.69515142907905952645E-15                            // 0x3ce8469898cc5170,
		M4PI = 1.273239544735162542821171882678754627704620361328125 // 4/pi
	)
	// special cases
	switch {
	case IsNaN(x) || IsInf(x, 0):
		return float32(math.NaN())
	}

	// make argument positive
	sign := false
	if x < 0 {
		x = -x
	}

	j := int32(x * M4PI) // integer part of x/(Pi/4), as integer for tests on the phase angle
	y := float32(j)      // integer part of x/(Pi/4), as float

	// map zeros to origin
	if j&1 == 1 {
		j++
		y++
	}
	j &= 7 // octant modulo 2Pi radians (360 degrees)
	if j > 3 {
		j -= 4
		sign = !sign
	}
	if j > 1 {
		sign = !sign
	}

	z := ((x - y*PI4A) - y*PI4B) - y*PI4C // Extended precision modular arithmetic
	zz := z * z
	if j == 1 || j == 2 {
		y = z + z*zz*((((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5])
	} else {
		y = 1.0 - 0.5*zz + zz*zz*((((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5])
	}
	if sign {
		y = -y
	}
	return y
}

func Sqrt(v float32) float32 {
	// TODO: Make it float32
	return float32(math.Sqrt(float64(v)))
}
