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

//IsUp returns true if resource is up
func (b *Breaker) IsUp() bool {
	isUp := atomic.LoadInt32(&b.Down) == 0
	if !isUp {
		b.resetIfDue()
	}
	return isUp
}

//FlagUp flags resource down
func (b *Breaker) F