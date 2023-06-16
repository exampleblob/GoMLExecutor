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
	fmt.Printf("m0 acquire proceeded since a goroutine waiting on the lock finished\n")

	fl := new(sync.WaitGroup)
	fl.Add(1)
	go func() {
		fmt.Printf("iL wait for acquire might need to wait for prior 2 goroutines\n")
		s.Acquire(bctx)
		defer s.Release()

		assert.Equal(t, s.r, int32(0), "main and latest goroutine should have locked semaphore")
		doneL = true

		fmt.Printf("iL done\n")
		fl.Done()
	}()

	fmt.Printf("m0 wait for i3 goroutine to run\n")
	// once waiting, the original 2 goroutines should've completed
	fl.Wait()

	fmt.Printf("m0 i3 goroutine has unblocked waitgroup\n")
	s.Release()

	for x := 0; x < numWaiters; x++ {
		assert.True(t, dones[x], fmt.Sprintf("inner bool %d", x+1))
	}

	assert.True(t, doneL, "last bool")
}

// not really a test
func TestOrdering(t *testing.T) {
	tests := 8

	s := NewSemaph(1)

	wgs := new(sync.WaitGroup)
	wgs.Add(tests)

	wgd := new(sync.WaitGroup)
	wgd.Add(tests)

	s.Acquire(context.Background())
	for i := 1; i < tests+1; i += 1 {
		ii := i
		go func() {
			wgs.Done()
			fmt.Printf("%d-acquire\n", ii)
			s.Acquire(context.Background())
			defer s.Release()
			fmt.Printf("%d-done\n", ii)
			wgd.Done()
		}()
	}

	wgs.Wait()
	s.Releas