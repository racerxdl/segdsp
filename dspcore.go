package main

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/spy2go"
	"log"
	"github.com/racerxdl/segdsp/demodcore"
)

const fifoSize = 1024 * 1024

var samplesFifo *fifo.Queue
var demodulator demodcore.DemodCore

func AddS16Fifo(data []spy2go.ComplexInt16) {
	samplesFifo.UnsafeLock()
	defer samplesFifo.UnsafeUnlock()
	for i := 0; i < len(data); i++ {
		if samplesFifo.UnsafeLen() >= fifoSize {
			log.Printf("FIFO Overflowing!!")
			break
		}

		var c = complex(float32(data[i].Imag) / 32768.0, float32(data[i].Real) / 32768.0)
		samplesFifo.UnsafeAdd(c)
	}
}