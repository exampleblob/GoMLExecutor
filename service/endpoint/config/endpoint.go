package config

import (
	"time"
)

const defaultMaxEvaluatorWait = 50

// Endpoint represents an endpoint
type Endpoint struct {
	Port           int
	ReadTimeoutMs  int           `json:",omitempty" yaml:",omitempty"`
	WriteTimeoutMs int           `json:",omitempty" yaml:",omitempty"`
	WriteTimeout   time.Duration `json:",omitempty" yaml:",omitempty"`
	MaxHeaderBytes int           `json:",omitempty" yaml:",omitempty"`

	// HTTP data buffer pool - used when reading a payload, for saving memory
	PoolMaxSize int `json:",omitempty" yaml:",omitempty"`
	BufferSize  int `json:",omitempty" yaml:",omitempty"`

	MaxEvaluatorWaitMs      int           `json:",omitempty" yaml:",omitempty"`
	MaxEvaluatorWait        time.Duration `json:",omitempty" yaml:",omitempty"`
	MaxEvaluatorConcurrency int64         `json:",omitempty" yaml:",omitempty"`
}

//Init init applied default settings
func (e *Endpoint) Init() {
	if e.Port == 0 {
		e.Port = 8080
	}
	if e.ReadTimeoutMs == 0 {
		e.ReadTimeoutMs = 5000
	}
	if e.WriteTimeou