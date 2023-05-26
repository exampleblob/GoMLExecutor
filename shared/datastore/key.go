package datastore

import (
	"fmt"
	"strconv"
	"time"

	aero "github.com/aerospike/aerospike-client-go"
	"github.com/viant/mly/shared/config"
)

// Key represents a datastore key
type Key struct {
	Namespace string
	Set       string
	Value     interface{}
	*aero.GenerationPolicy
	TimeToLive time.Duration
	L2         *Key
}

func (k *Key) AsString() string {
	switch value := k.Value.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)