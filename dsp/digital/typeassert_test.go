package digital

import (
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/dsp/digital/binarySlicer"
	"testing"
)

func TestComplexWorkers(t *testing.T) {
	var complexWorkersType = []interface{}{
		&ComplexClockRecovery{},
	}

	for _, v := range complexWorkersType {
		_, ok := v.(dsp.ComplexWorker)
		if !ok {
			t.Fatalf("Type %T does not implement ComplexWorker type!\n", v)
		}
	}
}

func TestFloat32Workers(t *testing.T) {
	var floatWorkersType = []interface{}{
		&FloatClockRecovery{},
	}

	for _, v := range floatWorkersType {
		_, ok := v.(dsp.Float32Worker)
		if !ok {
			t.Fatalf("Type %T does not implement Float32Worker type!\n", v)
		}
	}
}

func TestComplex2Float32Workers(t *testing.T) {
	var cfWorkersType = []interface{}{}

	for _, v := range cfWorkersType {
		_, ok := v.(dsp.Complex2Float32Worker)
		if !ok {
			t.Fatalf("Type %T does not implement Complex2Float32Worker type!\n", v)
		}
	}
}

func TestFloat322ByteWorkers(t *testing.T) {
	var fbWorkersType = []interface{}{
		&binarySlicer.Float2LevelSlicer{},
		&binarySlicer.Float4LevelSlicer{},
	}

	for _, v := range fbWorkersType {
		_, ok := v.(dsp.Float322ByteWorker)
		if !ok {
			t.Fatalf("Type %T does not implement Float322ByteWorkers type!\n", v)
		}
	}
}
