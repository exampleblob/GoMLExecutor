
package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/aerospike/aerospike-client-go/types"
	"github.com/viant/bintly"
	"github.com/viant/gmetric"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/config"
	"github.com/viant/mly/shared/datastore/client"
	"github.com/viant/mly/shared/stat"
	"github.com/viant/scache"
	"github.com/viant/toolbox"
)

type CacheStatus int

const (
	// CacheStatusFoundNoSuchKey we cache the status that we did not find a cache; this no-cache value has a shorter expiry
	CacheStatusFoundNoSuchKey = CacheStatus(iota)
	// CacheStatusNotFound no such key status
	CacheStatusNotFound
	// CacheStatusFound entry found status
	CacheStatusFound
)

// Service datastore service
type Service struct {
	config *config.Datastore
	mode   StoreMode

	useLocal bool
	cache    *scache.Cache
	l1Client *client.Service
	l2Client *client.Service

	readCounter  *gmetric.Operation
	writeCounter *gmetric.Operation
}

func (s *Service) Config() *config.Datastore {
	return s.config
}

func (s *Service) debug() bool {
	return s.config.Debug
}

func (s *Service) id() string {
	return s.config.ID
}

func (s *Service) Mode() StoreMode {
	return s.mode
}

func (s *Service) SetMode(mode StoreMode) {
	s.mode = mode
}

func (s *Service) Enabled() bool {
	if s.config == nil {
		return false
	}
	return !s.config.Disabled
}

func (s *Service) Key(key string) *Key {
	return NewKey(s.config, key)
}

// Put implements Storer.Put
func (s *Service) Put(ctx context.Context, key *Key, value Value, dictHash int) error {
	stats := stat.NewValues()

	if s.writeCounter != nil {
		onDone := s.writeCounter.Begin(time.Now())
		defer func() {
			onDone(time.Now(), *stats...)
		}()
	}

	// Add to local cache first
	if err := s.updateCache(key.AsString(), value, dictHash); err != nil {
		stats.Append(err)
		return err
	}

	if s.l1Client == nil || s.mode == ModeClient {
		return nil
	}

	storable := getStorable(value)
	bins, err := storable.Iterator().ToMap()
	if err != nil {
		stats.Append(err)
		return err
	}

	if dictHash != 0 {
		bins[common.HashBin] = dictHash
	}

	writeKey, _ := key.Key()

	isDebug := s.debug()
	if isDebug {
		log.Printf("[%s datastore put] l1 %+v bins %+v", s.id(), writeKey, bins)
	}

	wp := key.WritePolicy(0)
	wp.SendKey = true
	if s.l1Client != nil && !s.config.ReadOnly {
		if err = s.l1Client.Put(wp, writeKey, bins); err != nil {
			stats.Append(err)
			return err
		}

		stats.Append(stat.L1Write)

		if isDebug {
			log.Printf("[%s datastore put] l1 OK", s.id())
		}
	}

	if s.l2Client != nil && !s.config.L2.ReadOnly {
		k2Key, _ := key.L2.Key()
		err = s.l2Client.Put(wp, k2Key, bins)
		if err != nil {
			stats.Append(err)
		} else {
			stats.Append(stat.L2Write)
		}

		if isDebug {
			log.Printf("[%s datastore put] l2 err:%v", s.id(), err)
		}
	}
	return err
}

// GetInto implements Storer.GetInto
func (s *Service) GetInto(ctx context.Context, key *Key, storable Value) (dictHash int, err error) {
	return s.getInto(ctx, key, storable)
}

func (s *Service) getInto(ctx context.Context, key *Key, storable Value) (int, error) {
	stats := stat.NewValues()

	if s.readCounter != nil {
		onDone := s.readCounter.Begin(time.Now())
		defer func() {
			onDone(time.Now(), *stats...)
		}()
	}

	keyString := key.AsString()
	if s.useLocal {
		status, dictHash, err := s.readFromCache(keyString, storable, stats)
		if err != nil {
			return 0, err
		}
		switch status {
		case CacheStatusFoundNoSuchKey:
			return 0, types.ErrKeyNotFound
		case CacheStatusFound:
			return dictHash, nil
		}
	}

	if s.mode == ModeServer {
		// in server mode, cache hit rate would be low and expensive, thus skipping it
		return 0, types.ErrKeyNotFound
	}
	if s.l1Client == nil {
		return 0, types.ErrKeyNotFound
	}
	dictHash, err := s.getFromClient(ctx, key, storable, stats)
	if common.IsInvalidNode(err) {
		stats.Append(stat.Down)
		err = nil
		return 0, types.ErrKeyNotFound
	}

	if s.useLocal && key != nil {
		if err == nil {
			if storable != nil {
				err = s.updateCache(keyString, storable, dictHash)
			}
		} else if common.IsKeyNotFound(err) {
			if e := s.updateNotFound(keyString); e != nil {
				return 0, types.ErrKeyNotFound
			}
		}
	}
	return dictHash, err
}

func (s *Service) updateNotFound(keyString string) error {
	entry := &Entry{
		Key:      keyString,
		NotFound: true,
		Expiry:   time.Now().Add(s.config.RetryTime()),
	}
	data, err := bintly.Encode(entry)
	if err != nil {
		return err
	}
	return s.cache.Set(keyString, data)
}

func (s *Service) updateCache(keyString string, entryData EntryData, dictHash int) error {
	if entryData == nil {
		return fmt.Errorf("entry was nil")
	}