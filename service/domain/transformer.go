
package domain

import (
	"context"
	"fmt"

	"github.com/viant/gtly"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/common/storable"
)

// Transformer is an adapter module used when the output of the TensorFlow model wants to be modified server side.
// signature is the Signature of the relevant model, determined by request context.
// input is the request body unmarshalled.
// output is the TensorFlow SavedModel output.
type Transformer func(ctx context.Context, signature *Signature, input *gtly.Object, output interface{}) (common.Storable, error)

func appendPair(outputs []Output, v []interface{}, ii int, pairs []*kvPair) ([]*kvPair, error) {
	var outputValue interface{}

	for i, tensor := range v {
		switch t := tensor.(type) {
		case [][]string:
			outputValue = t[ii][0]
		case [][]float32:
			outputValue = t[ii][0]
		case [][]float64:
			outputValue = t[ii][0]
		case [][]int64:
			outputValue = t[ii][0]
		case [][]int32:
			outputValue = t[ii][0]
		case []string:
			outputValue = t[ii]
		case []float32:
			outputValue = t[ii]
		case []float64:
			outputValue = t[ii]
		case []int64:
			outputValue = t[ii]