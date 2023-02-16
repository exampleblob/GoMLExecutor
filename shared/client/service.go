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
		response.Status = common.StatusCached
		response.DictHash = dictHash
	}
	return cachedCount, nil
}

func (s *Service) readFromCacheInBatch(ctx context.Context, batchSize int, dataType reflect.Type, cachable Cachable, response *Response, cached []interface{}) (int, error) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(batchSize)
	var err error
	mux := sync.Mutex{}
	var cachedCount = 0
	for k := 0; k < batchSize; k++ {
		go func(index int) {
			defer waitGroup.Done()
			cacheEntry := reflect.New(dataType.Elem()).Interface()
			key := cachable.CacheKeyAt(index)
			has, dictHash, e := s.readFromCache(ctx, key, cacheEntry)
			mux.Lock()
			defer mux.Unlock()
			if e != nil {
				err = e
				return
			}
			if has {
				response.DictHash = dictHash
				cachable.FlagCacheHit(index)
				cached[index] = cacheEntry
				cachedCount++
			}
		}(k)
	}
	waitGroup.Wait()
	return cachedCount, err
}

func (s *Service) readFromCache(ctx context.Context, key string, target interface{}) (bool, int, error) {
	if s.datastore == nil || !s.datastore.Enabled() {
		return false, 0, nil
	}

	dataType := reflect.TypeOf(target)
	if dataType.Kind() != reflect.Ptr {
		return false, 0, fmt.Errorf("invalid response data type: expeted ptr but had: %T", target)
	}

	storeKey := s.datastore.Key(key)
	dictHash, err := s.datastore.GetInto(ctx, storeKey, target)
	if err == nil {
		if (!s.Config.DictHashValidation) || dictHash == 0 || dictHash == s.dictionary().hash {
			return true, dictHash, nil
		}
	}

	return false, 0, err
}

func (s *Service) releaseMessage(input interface{}) {
	releaser, ok := input.(Releaser)
	if ok {
		releaser.Release()
	}
}

func (s *Service) dictionary() *Dictionary {
	s.RWMutex.RLock()
	dict := s.dict
	s.RWMutex.RUnlock()
	return dict
}

func (s *Service) init(options []Option) error {
	for _, option := range options {
		option.Apply(s)
	}

	if s.gmetrics == nil {
		s.gmetrics = gmetric.New()
	}

	location := reflect.TypeOf(Service{}).PkgPath()
	s.counter = s.gmetrics.MultiOperationCounter(location, s.Model+"Client", s.Model+" client performance", time.Microsecond, time.Minute, 2, stat.NewStore())
	s.dictCounter = s.gmetrics.MultiOperationCounter(location, s.Model+"ClientDict", s.Model+" client dictionary performance", time.Microsecond, time.Minute, 1, stat.ErrorOnly())

	if s.ErrorHistory == nil {
		impl := misragries.NewMisraGries(20)
		s.ErrorHistory = mg.New(impl)
	}

	if s.Config.MaxRetry == 0 {
		s.Config.MaxRetry = 3
	}

	err := s.initHTTPClient()
	if err != nil {
		return err
	}

	if s.Config.Datastore == nil {
		if err := s.loadModelConfig(); err != nil {
			return err
		}
	}
	if s.dict == nil {
		if err := s.loadModelDictionary(); err != nil {
			return err
		}
	}
	if ds := s.Config.Da