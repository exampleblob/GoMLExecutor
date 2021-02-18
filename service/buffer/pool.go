package buffer

import "sync/atomic"

//Pool represents data pool
type Pool struct {
	channel     chan []byte
	poolMaxSize int32
	count       int32
	bufferSize  int
}

//Get returns bytes
func (p *Pool) Get() (result []byte) {
	select {
	case result = <-p.channel:
	default:
		result = make([]byte, p.bufferSize)
	}
	atomic.AddInt32(&p.count, -1)
	return result
}

//P