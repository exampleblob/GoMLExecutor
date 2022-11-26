
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/viant/mly/shared/common"
)

// Message represents the client-side perspective of the ML prediction.
// The JSON payload is built along the method calls; be sure to call (*Message).start() to set up the opening "{".
// TODO document how cache management is built into this type.
// There are 2 "modes" for building the message: single and batch modes.
// For single mode, the JSON object contents are written to Message.buf per method call.
// Single mode functions include:
//    (*Message).StringKey(string, string)
//    (*Message).IntKey(string, int)
//    (*Message).FloatKey(string, float32)
// Batch mode is initiated by called (*Message).SetBatchSize() to a value greater than 0.
// For batch mode, the JSON payload is generated when (*Message).end() is called.
// Batch mode functions include (the type name is plural):
//    (*Message).StringsKey(string, []string)
//    (*Message).IntsKey(string, []int)
//    (*Message).FloatsKey(string, []float32)
// There is no strict struct for request payload since some of the keys of the request are dynamically generated based on the model inputs.
// The resulting JSON will have property keys that are set based on the model, and two optional keys, "batch_size" and "cache_key".
// Depending on if single or batch mode, the property values will be scalars or arrays.
// See service.Request for server-side perspective.
// TODO separate out single and batch sized request to their respective calls endpoints; the abstracted polymorphism currently is more
// painful than convenient.
type (
	Message struct {
		mux  sync.RWMutex // locks pool
		pool *messages

		batchSize int

		buf   []byte // contains the JSON message as it is built
		index int

		buffer *bytes.Buffer // used to build cache key

		keys []string
		key  string // memoize join of keys

		// used to represent multi-row requests, with batchSize > 0
		keyLock   sync.Mutex // locks multiKeys
		multiKeys [][]string
		multiKey  []string // memoize keys
		transient []*transient

		cacheHits  []bool // in multi-row requests, indicates if cache has a value ofr the key
		dictionary *Dictionary
	}

	transient struct {
		name   string
		values interface{}
		kind   reflect.Kind
	}
)

// Strings is used to debug the current message.
func (m *Message) Strings() []string {
	fields := m.dictionary.Fields()
	if len(m.transient) == 0 {
		return nil
	}
	var result = make([]string, 0)
	for i := 0; i < m.batchSize; i++ {
		record := map[string]interface{}{}

		for _, trans := range m.transient {
			field, ok := fields[trans.name]
			if !ok {
				continue
			}
			var values []string
			switch actual := trans.values.(type) {
			case []string:
				values = actual
			}
			value := values[0]
			if i < len(values) {
				value = values[i]
			}
			record[field.Name] = value

		}
		if data, _ := json.Marshal(record); len(data) > 0 {
			result = append(result, string(data))
		}
	}

	return result
}

func (m *Message) CacheHit(index int) bool {
	if index < len(m.cacheHits) {
		return m.cacheHits[index]
	}
	return false
}

func (m *Message) SetBatchSize(batchSize int) {
	m.batchSize = batchSize
}

func (m *Message) BatchSize() int {
	return m.batchSize
}

// Size returns message size
func (m *Message) Size() int {
	return m.index
}

// start must be called before end()
func (m *Message) start() {
	m.appendByte('{')
}

// end completes the JSON payload
func (m *Message) end() error {
	if len(m.multiKeys) > 0 {
		if err := m.endInMultiKeyMode(); err != nil {
			return err
		}

		m.trim(',')
		m.appendString("}\n")
		return nil
	}

	if len(m.keys) > 0 {
		m.addCacheKey()
	}

	m.trim(',')
	m.appendString("}\n")
	return nil
}

// StringKey sets key/value pair
func (m *Message) StringKey(key, value string) {
	m.panicIfBatch("Strings")
	if key, index := m.dictionary.lookupString(key, value); index != unknownKeyField {
		m.keys[index] = key
	}
	m.appendByte('"')
	m.appendString(key)
	m.appendString(`":"`)
	m.appendString(value)
	m.appendString(`",`)
}

// panicIfBatch ensure that if multi keys are use no single message is allowed
func (m *Message) panicIfBatch(typeName string) {
	if m.batchSize > 0 {
		panic(fmt.Sprintf("use %vKey", typeName))
	}
}

// StringsKey sets key/values pair
func (m *Message) StringsKey(key string, values []string) {
	m.ensureMultiKeys(len(values))
	m.transient = append(m.transient, &transient{name: key, values: values, kind: reflect.String})
	var index fieldOffset
	var keyValue string
	for i, value := range values {
		if len(m.multiKeys[i]) == 0 {
			m.multiKeys[i] = make([]string, m.dictionary.inputSize())
		}

		if keyValue, index = m.dictionary.lookupString(key, value); index != unknownKeyField {
			m.multiKeys[i][index] = keyValue
		}
	}

	m.expandKeysIfNeeds(len(values), index, keyValue)

}

// IntsKey sets key/values pair
func (m *Message) IntsKey(key string, values []int) {
	m.ensureMultiKeys(len(values))
	m.transient = append(m.transient, &transient{name: key, values: values, kind: reflect.Int64})

	var index fieldOffset
	var intKeyValue int
	var keyValue string