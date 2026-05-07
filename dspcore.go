package main

import (
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/segdsp/demodcore"
)

type DSPPipeline struct {
	samplesFifo *fifo.Queue
	Demodulator demodcore.DemodCore
	buffer      []complex64
	cb          func(interface{})
	ready       chan struct{}
	done        chan struct{}
}

func NewDSPPipeline() *DSPPipeline {
	return &DSPPipeline{
		samplesFifo: fifo.NewQueue(),
		ready:       make(chan struct{}, 1),
		done:        make(chan struct{}),
	}
}

func (d *DSPPipeline) AddComplex(data []complex64) {
	d.samplesFifo.Add(data)
	select {
	case d.ready <- struct{}{}:
	default:
	}
}

func (d *DSPPipeline) SetCallback(cb func(interface{})) {
	d.cb = cb
}

func (d *DSPPipeline) Start() {
	go d.loop()
}

func (d *DSPPipeline) Stop() {
	close(d.ready)
	<-d.done
}

func (d *DSPPipeline) run() {
	length := d.samplesFifo.Len()

	if length == 0 {
		return
	}

	if d.Demodulator == nil || d.cb == nil {
		return
	}

	d.buffer = d.samplesFifo.Next().([]complex64)

	var out = d.Demodulator.Work(d.buffer)

	if out != nil {
		d.cb(out)
	}
}

func (d *DSPPipeline) loop() {
	defer close(d.done)
	for range d.ready {
		for d.samplesFifo.Len() > 0 {
			d.run()
		}
	}
}
