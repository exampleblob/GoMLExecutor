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
	if ds := s.Config.Datastore; ds != nil {
		ds.Init()
		if err = ds.Validate(); err != nil {
			return err
		}
	}

	if s.datastore == nil {
		err := s.initDatastore()
		return err
	}
	s.messages = NewMessages(s.dictionary)
	return nil
}

func (s *Service) initHTTPClient() error {
	host, _ := s.getHost()
	var tslConfig *tls.Config
	if host != nil && host.IsSecurePort() {
		cert, err := getCertPool()
		if err != nil {
			return fmt.Errorf("failed to create certificate: %v", err)
		}

		tslConfig = &tls.Config{
			RootCAs: cert,
		}
	}

	http2Transport := &http2.Transport{
		TLSClientConfig: tslConfig,
	}

	if host == nil || !host.IsSecurePort() {
		http2Transport.AllowHTTP = true
		http2Transport.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		}
	}

	s.httpClient.Transport = http2Transport
	return nil
}

func (s *Service) loadModelConfig() error {
	var err error
	host, err := s.getHost()
	if err != nil {
		return err
	}
	if s.Config.Datastore, err = s.discoverConfig(host, host.metaConfigURL(s.Model)); err != nil {
		return err
	}
	s.Config.updateCache()
	return nil
}

func (s *Service) loadModelDictionary() error {
	stats := stat.NewValues()

	onDone := s.dictCounter.Begin(time.Now())
	defer func() {
		onDone(time.Now(), stats.Values()...)
	}()

	host, err := s.getHost()
	if err != nil {
		stats.Append(err)
		return err
	}
	URL := host.metaDictionaryURL(s.Model)

	httpClient := s.getHTTPClient(host)
	response, err := httpClient.Get(URL)
	if err != nil {
		// no context errors supported
		stats.Append(err)
		return fmt.Errorf("failed to load Dictionary: %w", err)
	}

	if response.Body == nil {
		err = fmt.Errorf("unable to load dictioanry body was empty")
		stats.Append(err)
		return err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		stats.Append(err)
		return fmt.Errorf("failed to read body: %w", err)
	}

	dict := &common.Dictionary{}
	if err = json.Unmarshal(data, dict); err != nil {
		stats.Append(err)
		return fmt.Errorf("failed to unmarshal dict: %w", err)
	}

	s.RWMutex.Lock()
	s.dict = NewDictionary(dict, s.Datastore.Inputs)
	s.RWMutex.Unlock()

	s.messages = NewMessages(s.dictionary)
	return nil
}

func (s *Service) getHTTPClient(host *Host) *http.Client {
	httpClient := http.DefaultClient
	if host.IsSecurePort() {
		httpClient = &s.httpClient
	}
	return httpClient
}

func (s *Service) initDatastore() error {
	remoteCfg := s.Config.Datastore
	if remoteCfg == nil {
		return nil
	}

	if remoteCfg.Datastore.ID == "" {
		return nil
	}

	var stores = map[string]*datastore.Service{}
	var err error
	datastores := &sconfig.DatastoreList{
		Datastores:  []*sconfig.Datastore{&remoteCfg.Datastore},
		Connections: remoteCfg.Connections,
	}

	if stores, err = datastore.NewStores(datastores, s.gmetrics); err != nil {
		return err
	}

	s.datastore = stores[remoteCfg.ID]
	if len(remoteCfg.Fields) > 0 {
		if err := remoteCfg.FieldsDescriptor(remoteCfg.Fields); err != nil {
			return err
		}
		s.newStorable = func() common.Storable {
			return storable.New(remoteCfg.Fields)
		}
	}

	if s.datastore != nil {
		s.datastore.SetMode(datastore.ModeClient)
	}

	return nil
}

func (s *Service) Close() error {
	s.httpClient.CloseIdleConnections()
	if s.ErrorHistory != nil {
		s.ErrorHistory.Close()
	}

	return nil
}

// New creates new client.
func New(model string, hosts []*Host, options ...Option) (*Service, error) {
	for i := range hosts {
		hosts[i].Init()
	}
	aClient := &Service{
		Config: Config{
			Model: model,
			Hosts: hosts,
		},
	}
	err := aClient.init(options)
	return aClient, err
}

func (s *Service) discoverConfig(host *Host, URL string) (*config.Remote, error) {
	httpClient := s.getHTTPClient(host)

	response, err := httpClient.Get(URL)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(respo