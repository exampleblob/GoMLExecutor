package tfmodel

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"runtime/trace"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/mly/service/clienterr"
	"github.com/viant/mly/service/domain"
)

//Evaluator represents evaluator
type Evaluator struct {
	id string

	session *tf.Session
	fetches []tf.Output
	targets []*tf.Operation
	domain.Signature
}

func (e *Evaluator) feeds(feeds []interface{}) (map[tf.Output]*tf.Tensor, error) {
	var result = make(map[tf.Output]*tf.Tensor, len(feeds))
	for _, input := range e.Signature.Inputs {
		tensor, err := tf.NewTensor(feeds[input.Index])
		if err != nil {
			return nil, fmt.Errorf("failed to prepare feed: %v(%v), due to %w", input.Name, feeds[input.Index], err)
		}
		result[input.Placeholder] = tensor
	}
	return result, nil
}

//Evaluate evaluates model
func (e *Evaluator) Evaluate(params []interface{}) ([]interface{}, error) {
	ctx := context.Background()

	if TFSessionPanicDuration > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, TFSessionPanicDuration)
		defer cancel()
	}

	errc := make(ch