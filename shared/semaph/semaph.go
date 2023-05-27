package semaph

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Semaph struct {
	l   sync.Mutex // 