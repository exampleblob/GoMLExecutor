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
	for _, inpu