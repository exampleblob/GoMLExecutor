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
func (h *Host) evalURL(model string) string {
	if h.IsSecurePort() {
		return "https://" + h.Name + ":" + strconv.Itoa(h.Port) + fmt.Sprintf(common.ModelURI, model)
	}
	return "http://" + h.Name + ":" + strconv.Itoa(h.Port) + fmt.Sprintf(common.ModelURI, model)
}

//URL returns meta config model eval URL
func (h *Host) metaConfigURL(model string) string {
	if h.IsSecu