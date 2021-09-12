package service

import (
	"compress/gzip"
	"context"
	sjson "encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"runtime/trace"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/gmetric"
	"github.com/viant/gtly"
	"github.com/viant/mly/service/clienterr"
	"github.com/viant/mly/service/config"
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/service/files"
	"github.com/viant/mly/service/layers"
	"github.com/viant/mly/service/request"
	"github.com/viant/mly/service/stat"
	"github.com/viant/mly/service/stream"
	"github.com/viant/mly/service/tfmodel"
	"github.com/viant/mly/service/transform"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/common/storable"
	"github.com/viant/mly/shared/datastore"
	sstat "github.com/viant/mly/shared/stat"
	"github.com/viant/xunsafe"
	"golang.org/x/sync/semaphore"
	"gopkg.in/yaml.v3"
)

type Service struct {
	config *config.Model
	closed int32

	maxEvaluatorWait time.Duration

	// TODO how does this interact with Service.inputs
	inputProvider *gtly.Provider

	// reload
	ReloadOK int32
	fs       afs.Service
	mux      sync.RWMutex

	// model
	sema      *semaphore.Weighted // prevents potentially explosive thread generation due to concurrent requests
	evaluator *tfmodel.Evaluator

	// model io
	signature *domain.Signature
	inputs    map[string]*domain.Input

	// caching
	useDatastore bool
	dictionary   *common.Dictionary
	datastore    datastore.Storer

	// outputs
	transformer domain.Transformer
	newStorable func() common.Storable

	// metrics
	serviceMetric   *gmetric.Operation
	evaluatorMetric *gmetric.Operation

	// logging
	stream *stream.Service
}

func (s *Service) Clos