package storable

import (
	"reflect"
	"strings"
	"sync"
)

var _reflect = &reflectCache{cache: make(map[string]*reflectStruct)}

type refle