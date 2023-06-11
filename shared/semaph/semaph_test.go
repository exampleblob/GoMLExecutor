package semaph

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestCore(t *testing.T) {
	semaSize := 2
	s := NewSemaph(int32(semaSize))

	numWaiters := 3

	waitAcq := new(sync.WaitGroup)
	waitAcq.Add(numWaiters)

	waitAcqd := new(sync.WaitGroup)
	waitAcqd.Add(semaSize)

	ol := new(sync.Mutex)
	ol.Lock()

	dones := make([]bool, numWaiters)
	var doneL bool

	bctx := context.Background()

	// since semaphore size is 2 and num workers is 3, we should have 1 worker block
	for x := 1; x <= numWaiters; x++ {
		go func(id int, marker *bool) {
			waitAcq.Done()
			fmt.Printf("i%d acquiring\n", id)
			s.Acquire(bctx)
			defer s.Release()
			waitAcqd.Done()

			fmt.Printf("i%d wait for m0 lock\n", id)
			ol.Lock()
			defer ol.Unlock()
			*marker = true

			fmt.Printf("i%d done\n", id)
		}(x, &dones[x-1])
	}

	fmt.Printf("m0 wait for semaphores to be acquired by goroutines\n")
	waitAcq.Wait()
	fmt.Printf("m0 goroutines should have acquired\n")

	waitAcqd.Wait()

	// prevent error from calling wg.Done() too many times
	waitAcqd.Add(numWaiters - semaSize)

	assert.Equal(t, s.r, int32(0), "the sempahore should not be available")

	fmt.Printf("m0 unlock as enough goroutines should have started\n")
	ol.Unlock()
	fmt.Printf("m0 unlocked\n")

	s.Acquire(bctx)
	fmt.Printf("m0 acquire proceeded since a goroutine waiting on the lock finish