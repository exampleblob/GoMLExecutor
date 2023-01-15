package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cyningsun/heavy-hitters/misragries"
	"github.com/francoispqt/gojay"
	"github.com/viant/gmetric"
	"github.com/viant/mly/shared/client/config"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/common/storable"
	sconfig "github.com/viant/mly/shared/config"
	"github.com/viant/mly/shared/datastore"
	"github.com/viant/mly/shared/stat"
	"github.com/viant/mly/shared/tracker"
	"github.com/viant/mly/shared/tracker/mg"
	"github.com/viant/xunsafe"
	"golang.org/x/net/http2"
)

//Service represent mly client
type Service struct {
	Config
	hostIndex   int64
	httpClient  http.Client
	newStorable func() common.Storable

	messages Messages
	poolErr  error

	sync.RWMutex
	dict               *Dictionary
	dictRefreshPending int32
	datastore          datastore.Storer

	// container for Datastore gmetric objects
	gmetrics *gmetric.Service

	counter     *gmetric.Operation
	dictCounter *gmetric.Operation

	ErrorHistory tracker.Tracker
}

// NewMessage returns a new message
func (s *Service) NewMessage() *Message {
	message := s.messages.Borrow()
	message.start()
	return message
}

// Run run model prediction
func (s *Service) Run(ctx context.Context, input interface{}, response *Response) error {
	onDone := s.counter.Begin(time.Now())
	stats := stat.NewValues()
	defer func() {
		onDone(time.Now(), *stats...)
		s.releaseMessage(input)
	}()

	if response.Data == nil {
		return fmt.Errorf("response data was empty - aborting request")
	}

	cachable, isCachable := input.(Cachable)
	var err error
	var cachedCount int
	var cached []interface{}

	batchSize := 0
	if isCachable {
		batchSize = cachable.BatchSize()

		cachedCount, err = s.loadFromCache(ctx, &cached, batchSize, response, cachable)
		if err != nil {
			return err
		}
	}

	isDebug := s.Config.Debug

	var modelName string
	if isDebug {
		modelName = s.Config.Model
		s.reportBatch(cachedCount, cached)
	}

	if (batchSize > 0 && cachedCount == batchSize) || (batchSize == 0 && cachedCount > 0) {
		response.Status = common.StatusCached
		return s.handleResponse(ctx, response.Data, cached, cachable)
	}

	data, err := Marshal(input, modelName)
	if err != nil {
		return err
	}

	stats.Append(sta