
package stream

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/viant/afs"
	"github.com/viant/gmetric"
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/stat"
	"github.com/viant/tapper/config"
	"github.com/viant/tapper/io"
	tlog "github.com/viant/tapper/log"
	"github.com/viant/tapper/msg"
	"github.com/viant/tapper/msg/json"
)

type dictProvider func() *common.Dictionary
type outputsProvider func() []domain.Output

// Service is used to log request inputs to model outputs without an output
// transformer, in JSON format.
//
// The input values will be directly inlined into the resulting JSON.
// The outputs will be provided as properties in the resulting JSON, with
// the keys as the output Tensor names.
//
// If the dimensions of the output from the model are [1, numOutputs, 1] (single
// request), the value in the JSON object will be a scalar.
// If the dimensions of the output from the model are [batchSize, numOutputs, 1],
// (batch request), the value in the JSON object will be a list of scalars of
// length batchSize.
// If the dimensions of the output from the model are [1, numOutputs, outDims],
// (single request), the value of the JSON object will be a list of scalars of