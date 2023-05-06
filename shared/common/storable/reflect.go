package storable

import (
	"reflect"
	"strings"
	"sync"
)

var _reflect = &reflectCache{cache: make(map[string]*reflectStruct)}

type reflectCache struct {
	cache map[string]*reflectStruct
	mux   sync.RWMutex
}

func (c *reflectCache) lookup(aType reflect.Type) *reflectStruct {
	c.mux.RLock()
	result, ok := c.cache[aType.Nam