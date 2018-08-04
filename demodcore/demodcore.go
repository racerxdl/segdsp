package demodcore


type DemodCore interface {
	Work(data []complex64) interface{}
	GetDemodParams() interface{}
}