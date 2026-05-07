package recorders

const RecFile = "file"
const RecNone = "none"

type RecorderConfig struct {
	BaseFilename string
	Metadata     interface{}
	AudioAsWav   bool
}

type BaseRecorder interface {
	Open(config RecorderConfig) bool
	Close() bool
	WriteIQ(data []complex64)
	WriteAudio(data []float32)
	WriteData(data []byte)
}
