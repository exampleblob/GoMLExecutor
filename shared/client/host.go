package client

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/viant/mly/shared/circut"
	"github.com/viant/mly/shared/common"
)

var requestTimeout = 50 * time.Millisecond

//Host represents endpoint host
type Host struct {
	Name string
	Port int
	mux  sync.RWMutex
	*circut.Breaker
}

//IsSecurePort() returns true if secure port
func (h *Host) IsSecurePort() bool {
	return h.Port%1000 == 443
}

//URL returns model eval URL
func (h *Host) evalURL(model string) st