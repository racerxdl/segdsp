package dsp

import (
	"testing"
)

func TestComplexWorkers(t *testing.T) {
	var complexWorkersType = []interface{}{
		&CTFirFilter{},
		&DCFilter{},
		&Decimator{},
		&FeedForwardAGC{},
		&FirFilter{},
		&FrequencyTranslator{},
		&Interpolator{},
		&RationalResampler{},
		&Squelch{},
		&Rotator{},
	}

	for _, v := range complexWorkersType {
		_, ok := v.(ComplexWorker)
		if !ok {
			t.Fatalf("Type %T does not implement ComplexWorker type!\n", v)
		}
	}
}

func TestFloat32Workers(t *testing.T) {
	var floatWorkersType = []interface{}{
		&FloatFirFilter{},
		&FloatDecimator{},
		&FloatInterpolator{},
		&FloatRationalResampler{},
		&FloatResampler{},
	}

	for _, v := range floatWorkersType {
		_, ok := v.(Float32Worker)
		if !ok {
			t.Fatalf("Type %T does not implement Float32Worker type!\n", v)
		}
	}
}

func TestComplex2Float32Workers(t *testing.T) {
	var cfWorkersType = []interface{}{
		&QuadDemod{},
		&Complex2Magnitude{},
	}

	for _, v := range cfWorkersType {
		_, ok := v.(Complex2Float32Worker)
		if !ok {
			t.Fatalf("Type %T does not implement Float32Worker type!\n", v)
		}
	}
}
