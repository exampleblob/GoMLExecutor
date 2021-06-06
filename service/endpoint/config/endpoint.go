package config

import (
	"time"
)

const defaultMaxEvaluatorWait = 50

// Endpoint represents an endpoint
type Endpoint struct {
	Port           int
	ReadTimeoutMs  int           `json:",omitempty" yam