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
	targets []*tf.Operat