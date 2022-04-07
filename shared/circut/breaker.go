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
		return
	}
	b.resetTime = time.Now().Add(b.resetDuration)
	b.resetDuration = time.Duration(float32(b.resetDuration) * 1.5)
	b.mux.Unlock()
	go b.prober.Probe()
}

//FlagDown flags connection Down
func (b *Breaker) FlagDown() {
	down := atomic.LoadInt32(&b.Down)
	if down == 1 {
		return
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.Down == 1 {
		return
	}
	b.Down = 1
	b.resetTime = time.Now().Add(b.resetDurat