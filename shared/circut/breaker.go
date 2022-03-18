package circut

import (
	"sync"
	"sync/atomic"
	"time"
)

//Breaker represents circut breaker
type Breaker struct {
	prober               Prober
	Down                 int32
	mux                  sync.RWMutex
	resetTime            time.Time
	resetDuration        time.Duration
	initialResetDuration time.Duration
}

//IsUp returns true i