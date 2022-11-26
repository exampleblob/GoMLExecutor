
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