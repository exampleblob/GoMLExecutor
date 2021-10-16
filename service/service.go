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

func (s *Service) Close() error {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return fmt.Errorf("already closed")
	}

	if s.evaluator == nil {
		return nil
	}

	return s.evaluator.Close()
}

func (s *Service) Config() *config.Model {
	return s.config
}

func (s *Service) Dictionary() *common.Dictionary {
	return s.dictionary
}

func (s *Service) Do(ctx context.Context, request *request.Request, response *Response) error {
	err := s.do(ctx, request, response)
	if err != nil {
		response.Error = err.Error()
		response.Status = common.StatusError
		return err
	}

	return nil
}

func (s *Service) do(ctx context.Context, request *request.Request, response *Response) error {
	startTime := time.Now()
	onDone := s.serviceMetric.Begin(startTime)
	onPendingDone := incrementPending(s.serviceMetric, startTime)
	stats := sstat.NewValues()
	defer func() {
		onDone(time.Now(), stats.Values()...)
		onPendingDone()
	}()

	err := request.Validate()
	if s.config.Debug && err != nil {
		log.Printf("[%v do] validation error: %v\n", s.config.ID, err)
	}

	if err != nil {
		// only captures missing fields
		stats.Append(stat.Invalid)
		return clienterr.Wrap(fmt.Errorf("%w, body: %s", err, request.Body))
	}

	if err != nil {
		stats.AppendError(err)
		log.Printf("[%v do] limiter error:(%+v) request:(%+v)", s.config.ID, err, request)
		return err
	}

	cancel := func() {}
	if s.maxEvaluatorWait > 0 {
		// this is here due to how the semaphore operates
		ctx, cancel = context.WithTimeout(ctx, s.maxEvaluatorWait)
	}

	tensorValues, err := s.evaluate(ctx, request)
	cancel()

	if err != nil {
		// we waited or there was an issue with evaluation; in either case
		// the prediction never finished so there is nothing left to clean up
		stats.AppendError(err)
		log.Printf("[%v do] eval error:(%+v) request:(%+v)", s.config.ID, err, request)
		return err
	}

	stats.Append(stat.Evaluate)
	return s.buildResponse(ctx, request, response, tensorValues)
}

func (s *Service) transformOutput(ctx context.Context, request *request.Request, output interface{}) (common.Storable, error) {
	inputIndex := inputIndex(output)
	inputObject := request.Input.ObjectAt(s.inputProvider, inputIndex)

	transformed, err := s.transformer(ctx, s.signature, inputObject, output)
	if err != nil {
		return nil, fmt.Errorf("failed to transform: %v, %w", s.config.ID, err)
	}

	if s.useDatastore {
		dictHash := s.Dictionary().Hash
		cacheKey := request.Input.KeyAt(inputIndex)

		isDebug := s.config.Debug
		key := s.datastore.Key(cacheKey)

		go func() {
			err := s.datastore.Put(ctx, key, transformed, dictHash)
			if err != nil {
				log.Printf("[%s trout] put error:%v", s.config.ID, err)
			}

			if isDebug {
				log.Printf("[%s trout] put:\"%s\" dictHash:%d ok", s.config.ID, cacheKey, dictHash)
			}
		}()
	}

	return transformed, nil
}

func inputIndex(output interface{}) int {
	inputIndex := 0
	if out, ok := output.(*shared.Output); ok {
		inputIndex = out.InputIndex
	}
	return inputIndex
}

func (s *Service) buildResponse(ctx context.Context, request *request.Request, response *Response, tensorValues []interface{}) error {
	if s.dictionary != nil {
		response.DictHash = s.dictionary.Hash
	}

	// TODO change with understanding that batched / multi-request always operates on the first dimension
	if !request.Input.BatchMode() {
		var err error
		// single input
		response.Data, err = s.transformOutput(ctx, request, tensorValues)
		return err
	}

	output := &shared.Output{Values: tensorValues}
	// index 0 call to get data type
	transformed, err := s.transformOutput(ctx, request, output)
	if err != nil {
		return err
	}

	sliceType := reflect.SliceOf(reflect.TypeOf(transformed))
	// xSlice = make([]`sliceType`, request.Input.BatchSize)
	sliceValue := reflect.MakeSlice(sliceType, 0, request.Input.BatchSize)
	slicePtr := xunsafe.ValuePointer(&sliceValue)
	xSlice := xunsafe.NewSlice(sliceType)

	response.xSlice = xSlice
	response.sliceLen = request.Input.BatchSize

	// xSlice = append(xSlice, transformed)
	appender := xSlice.Appender(slicePtr)
	appender.Append(transformed)

	// index 1 - end calls
	for i := 1; i < request.Input.BatchSize; i++ {
		output.InputIndex = i
		if transformed, err = s.transformOutput(ctx, request, output); err != nil {
			return err
		}
		appender.Append(transformed)
	}

	response.Data = sliceValue.Interface()
	return nil
}

func incrementPending(metric *gmetric.Operation, startTime time.Time) func() {
	return incrementThenDecrement(metric, startTime, stat.Pending)
}

func incrementThenDecrement(metric *gmetric.Operation, start time.Time, statName string) func() {
	metric.IncrementValue(statName)

	index := metric.Index(start)
	recentCounter := metric.Recent[index]
	recentCounter.IncrementValue(statName)

	return func() {
		metric.DecrementValue(statName)
		recentCounter := metric.Recent[index]
		recentCounter.DecrementValue(statName)
	}

}

func (s *Service) evaluate(ctx context.Context, request *request.Request) ([]interface{}, error) {
	trr := trace.StartRegion(ctx, "Semaphore.Acquire")
	err := s.sema.Acquire(ctx, 1)
	trr.End()
	if err != nil {
		return nil, err
	}
	// even if canceled/deadline exceeded, we'r