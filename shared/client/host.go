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
	if h.IsSecurePort() {
		return "https://" + h.Name + ":" + strconv.Itoa(h.Port) + fmt.Sprintf(common.MetaConfigURI, model)
	}
	return "http://" + h.Name + ":" + strconv.Itoa(h.Port) + fmt.Sprintf(common.MetaConfigURI, model)
}

//URL returns meta config model eval URL
func (h *Host) metaDictionaryURL(model string) string {
	if h.IsSecurePort() {
		return "https://" + h.Name + ":" + strconv.Itoa(h.Port) + fmt.Sprintf(common.MetaDictionaryURI, model)
	}
	return "http://" + h.Name + ":" + strconv.Itoa(h.Port) + fmt.Sprintf(common.MetaDictionaryURI, model)
}

func (h *Host) Probe() {
	if h.isConnectionUp() {
		h.FlagUp()
		return
	}
}

func (h *Host) isConnectionUp() bool {
	port := h.Port
	connection, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", h.Name, port), requestTimeout)
	if err != nil {
		return false
	}
	defer connection.Close()
	return true
}

func (h *Host) Init() {
	if h.Breaker == nil {
		h