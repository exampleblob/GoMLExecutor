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
			fmt.Printf