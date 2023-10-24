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

func Run(ts TestStruct, maxDos int32, testCases 