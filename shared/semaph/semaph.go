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
	WaitCanceled uint64
}

func NewSemaph(max int32) *Semaph {
	s := new(Semaph)
	s.r = max
	s.max = max

	s.l = sync.Mutex{}
	s.c = sync.NewCond(&s.l)

	return s
}

func (s *Semaph) Internals() (int32, Stats) {
	r