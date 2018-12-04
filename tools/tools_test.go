package tools

import (
	"math"
	"math/rand"
	"testing"
)

const testRuns = 1 << 12

func DoTestAgainstStdGo(name string, f32 func(float32) float32, stdGo func(float64) float64, t *testing.T) {
	for i := 0; i < testRuns; i++ {
		var v = rand.Float32()*2 - 1
		got := f32(v)
		expected := stdGo(float64(v))

		if !AlmostFloatEqual(got, float32(expected)) {
			t.Errorf("Float32 %s wrong. Expected (%f) got (%f)", name, expected, got)
		}
	}
}

func TestAbs(t *testing.T) {
	DoTestAgainstStdGo("Abs", Abs, math.Abs, t)
}

func TestFloor(t *testing.T) {
	DoTestAgainstStdGo("Floor", Floor, math.Floor, t)
}

func TestAtan(t *testing.T) {
	DoTestAgainstStdGo("Atan", Atan, math.Atan, t)
}

func TestCopysign(t *testing.T) {
	for i := 0; i < testRuns; i++ {
		var a = rand.Float32()*2 - 1
		var b = rand.Float32()*2 - 1
		got := Copysign(a, b)
		expected := math.Copysign(float64(a), float64(b))

		if !AlmostFloatEqual(got, float32(expected)) {
			t.Errorf("Float32 Copysign wrong. Expected %f got %f", expected, got)
		}
	}
}

func TestHypot(t *testing.T) {
	for i := 0; i < testRuns; i++ {
		var a = rand.Float32()*2 - 1
		var b = rand.Float32()*2 - 1
		got := Hypot(a, b)
		expected := math.Hypot(float64(a), float64(b))

		if !AlmostFloatEqual(got, float32(expected)) {
			t.Errorf("Float32 Hypot wrong. Expected %f got %f", expected, got)
		}
	}
}

func TestModf(t *testing.T) {
	for i := 0; i < testRuns; i++ {
		var v = rand.Float32()*2 - 1
		gint, gfrac := Modf(v)
		eint, efrac := math.Modf(float64(v))

		if !AlmostFloatEqual(gint, float32(eint)) || !AlmostFloatEqual(gfrac, float32(efrac)) {
			t.Errorf("Float32 Modf wrong. Expected (%f, %f) got (%f, %f)", gint, gfrac, eint, efrac)
		}
	}
}

func TestAlmostFloatEqual(t *testing.T) {
	for i := 0; i < testRuns; i++ {
		// Not sure if this is an adequate test
		var v = rand.Float32()*2 - 1
		var variation = ((rand.Float32()*2 - 1) * epsilon) / 2

		if variation > epsilon || variation < -epsilon {
			t.Errorf("Test is wrong %f, %f", variation, epsilon)
		}

		if !AlmostFloatEqual(v, v+variation) || !AlmostFloatEqual(v, v-variation) {
			t.Errorf("Float32 Equal was wrong.\n Epsilon: [%f]\n AlmostFloatEqual(v, v + variation) => (%f, %f) = %v\n AlmostFloatEqual(v, v - variation) => (%f, %f) = %v", epsilon, v, v+variation, AlmostFloatEqual(v, v+variation), v, v-variation, AlmostFloatEqual(v, v-variation))
		}
	}
}

func TestSignbit(t *testing.T) {
	for i := 0; i < testRuns; i++ {
		var v = rand.Float32()*2 - 1
		got := Signbit(v)
		expected := v < 0

		if got != expected {
			t.Errorf("Float32 SignBit wrong. Expected %v got %v for %v ", expected, got, v)
		}
	}
}

func TestIsNaN(t *testing.T) {
	for i := 0; i < testRuns; i++ {
		var v = rand.Float32()*2 - 1
		got := IsNaN(v)
		expected := math.IsNaN(float64(v))

		if got != expected {
			t.Errorf("Float32 IsNaN wrong. Expected %v got %v for %v ", expected, got, v)
		}
	}
}
