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
	mux  sync.R