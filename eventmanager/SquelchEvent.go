package eventmanager

const EvSquelchOn = "squelchOnEvent"
const EvSquelchOff = "squelchOffEvent"

type SquelchEventData struct {
	Threshold float32
	AvgValue  float32
}
