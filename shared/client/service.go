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

	stats.Append(stat.NoSuchKey)
	if isDebug {
		fmt.Printf("[%s] request: %s\n", modelName, strings.Trim(string(data), " \n"))
	}

	body, err := s.postRequest(ctx, data)
	if isDebug {
		fmt.Printf("[%s] response.Body:%s\n", modelName, body)
		fmt.Printf("[%s] error:%s\n", modelName, err)
	}

	if err != nil {
		stats.AppendError(err)

		if ctx.Err() == nil && s.ErrorHistory != nil {
			go s.ErrorHistory.AddBytes([]byte(err.Error()))
		}

		return err
	}

	err = gojay.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: '%s'; due to %w", body, err)
	}

	if isDebug {
		fmt.Printf("[%v] response.Data: %s, %v\n", modelName, response.Data, err)
	}

	if response.Status != common.StatusOK {
		// TODO is this correct?
		return nil
	}

	if err = s.handleResponse(ctx, response.Data, cached, cachable); err != nil {
		return fmt.Errorf("failed to handle resp: %w", err)
	}

	s.updatedCache(ctx, response.Data, cachable, s.dict.hash)
	s.assertDictHash(response)

	return nil
}

func (s *Service) loadFromCache(ctx context.Context, cached *[]interface{}, batchSize int, response *Response, cachable Cachable) (int, error) {
	*cached = make([]interface{}, batchSize)
	dataType, err := response.DataItemType()
	if err != nil {
		return 0, err
	}

	if batchSize > 0 {
		cachedCount, err := s.readFromCacheInBatch(ctx, batchSize, dataType, cachable, response, *cached)
		if err != nil && !common.IsTransientError(err) {
			log.Printf("cache error: %v", err)
		}

		return cachedCount, nil
	}

	key := cachable.CacheKey()
	has, dictHash, err := s.readFromCache(ctx, key, response.Data)
	if err != nil && !common.IsTransientError(err) {
		log.Printf("cache error: %v", err)
	}
	cachedCount := 0
	if has {
		cachedCount = 1
		response.S