package smasher

import (
	"context"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/viant/mly/shared/semaph"
)

type (
	Server interface {
		Stats() string
	}

	Client interface {
		Do() error

		Sent() uint64
	}

	TestStruct struct {
		Server func() (Server, error)
		Client func() (Client, error)
	}
)

func Run(ts TestStruct, maxDos int32, testCases int, statDur time.Duration) error {
	srv, err := ts.Server()
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(testCases)

	cli, err := ts.Client()
	if err != nil {