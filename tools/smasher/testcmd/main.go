
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/viant/mly/tools/smasher"
)

type cl struct {
	cli  *http.Client
	URL  string
	sent uint64
}

func (c *cl) Do() error {
	_, err := c.cli.Get(c.URL)
	atomic.AddUint64(&c.sent, 1)
	return err
}