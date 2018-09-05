package demodcore

import "github.com/racerxdl/segdsp/eventmanager"

type DemodCore interface {
	Work(data []complex64) interface{}
	GetDemodParams() interface{}
	SetEventManager(ev *eventmanager.EventManager)
	GetLevel() float32
	IsMuted() bool
}