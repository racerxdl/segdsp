package recorders

const RecFile = "file"
const RecNone = "none"

type BaseRecorder interface {
	Open(params []interface{}) bool
	Close() bool
	WriteIQ(data []complex64)
	WriteAudio(data []float32)
	WriteData(data []byte)
}
