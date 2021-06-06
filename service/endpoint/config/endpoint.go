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
	MaxHeaderBytes int           `json:",omitempty" yaml:",omi