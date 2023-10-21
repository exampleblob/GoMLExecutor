
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/viant/mly/service/endpoint"
	"github.com/viant/mly/shared/client"
	"github.com/viant/mly/tools/smasher"
)

type mlyS struct {
	cli        *http.Client
	metricPath string
	ma         uint64
	l          sync.Mutex
}

type stats struct {
	Count    int
	Counters []counter
}

type counter struct {
	Value string
	Count int
}

func (s *mlyS) Stats() string {
	stp := s.getStats("Perf")
	ste := s.getStats("Eval")

	var pErr, eErr string
	if len(stp.Counters) > 0 {
		pErr = "noPerf"
	}

	if len(ste.Counters) > 0 {
		eErr = "noEval"
	}

	var perfPending, eval int
	for _, c := range stp.Counters {
		if c.Value == "pending" {
			perfPending = c.Count
		}
