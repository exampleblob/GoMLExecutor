package buffer

import "sync/atomic"

//Pool represents data pool
type Pool struct {
	channel     chan []byte
	poolMaxSize int32
	count       int32
	buff