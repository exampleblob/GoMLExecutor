package client

import (
	"log"

	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
)

type fieldOffset int

const (
	// oov = out of vocabulary
	oovString = "[UNK]"
	oovInt    = 0

	defaultPrec = 10

	unknownKeyField = fieldOffset(-1)
)

// Dictionary helps identify any out-of-vocabulary input values for reducing the cache space - this enables us to leverage any
// dimensionality reduction within the model to optimize wall-clock performance. This is primarily useful for categorical inputs
// as well as any continous inputs with an acceptable quantization.
type Dictionary struct {
	hash     int
	registry map[string]*entry
	inputs   map[string]*shared.Field
}

func (d *Dictionary) KeysLen() int {
	return len(d.inputs)
}

func (d *Dictionary) inputSize() int {
	return len(d.inputs)
}

func (d *Dictionary) size() int {
	return len(d.registry)
}

// TODO refactor, this has a singular use case
func (d *Dictionary) Fields() map[string]*shared.Field {
	return d.inputs
}

func (d *Dict