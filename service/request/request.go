
package request

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/francoispqt/gojay"
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/transfer"
)

var exists struct{} = struct{}{}

// Request represents the server-side post-processed information about a request.
// There is no strict struct for request payload since some of the keys of the request are dynamically generated based on the model inputs.
// See shared/client.Message for client-side perspective.
type Request struct {
	Body     []byte // usually the POST JSON content
	Feeds    []interface{}
	inputs   map[string]*domain.Input
	supplied map[string]struct{}
	Input    *transfer.Input
}

func NewRequest(keyLen int, inputs map[string]*domain.Input) *Request {
	return &Request{
		inputs:   inputs,
		Feeds:    make([]interface{}, keyLen),
		Input:    &transfer.Input{},
		supplied: make(map[string]struct{}, keyLen),
	}
}

// Put is used when constructing a request NOT using gojay.
func (r *Request) Put(key string, value string) error {
	if input, ok := r.inputs[key]; ok {