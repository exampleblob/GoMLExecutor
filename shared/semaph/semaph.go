package semaph

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Semaph struct {
	l   sync.Mutex // locks for modifying r
	c   *sync.Cond
	r   int32 // i.e. remaining tickets
	max int32

	stats Stats
}

type Stats struct {
	Acquired,
	Waited,
	Canceled,
	CanceledDone,
	WaitDone,
	WaitCancele