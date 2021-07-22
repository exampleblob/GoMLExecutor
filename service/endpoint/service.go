
package endpoint

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/viant/gmetric"
	srvConfig "github.com/viant/mly/service/config"
	"github.com/viant/mly/service/endpoint/checker"
	"github.com/viant/mly/service/endpoint/health"
	"github.com/viant/mly/service/endpoint/prometheus"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/client"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/datastore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const healthURI = "/v1/api/health"
const statsURI = "/v1/api/stats"

// Service is the primary container for all HTTP based services.
type Service struct {
	server *http.Server
	config *Config
}

// Deprecated: use Listen and Serve separately
func (s *Service) ListenAndServe() error {
	ln, err := s.Listen()
	if err != nil {
		return err
	}

	return s.Serve(ln)
}
