
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