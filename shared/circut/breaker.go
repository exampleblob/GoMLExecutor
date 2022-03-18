package circut

import (
	"sync"
	"sync/atomic"
	"time"
)

//Breaker represents circut breaker
type Breaker struct {
	prober               Prober
	D