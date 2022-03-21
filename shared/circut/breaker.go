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
func (b *Breaker) FlagUp() {
	b.mux.Lock()
	b.Down = 0
	b.mux.Unlock()
	b.resetDuration = b.initialResetDuration
}

//resetIfDue reset connection onDisconnect status if reset time is due,
func (b *Breaker) resetIfDue() {
	b.mux.RLock()
	dueTime := time.Now().After(b.resetTime)
	b.mux.RUnlock()
	if !dueTime {
		return
	}
	b.mux.Lock()
	dueTime = time.Now().After(b.resetTime)
	if !dueTime {
		b.mux.Unlock()
		ret