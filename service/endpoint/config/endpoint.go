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

	MaxEvaluatorWaitMs      int 