package endpoint

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/viant/afs"
	"github.com/viant/gmetric"
	"github.com/viant/mly/service"
	"github.com/viant/mly/service/buffer"
	"github.com/viant/mly/service/config"
	serviceConfig "github.com/viant/mly/service/config"
	"github.com/viant/mly/service/endpoint/meta"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/datastore"
	"golang.org/x/sync/semaphore"
)

type Hook interface {
	Hook(*config.Model, *service.Service)
}

func Build(mux *http.ServeMux, config *Config, datasto