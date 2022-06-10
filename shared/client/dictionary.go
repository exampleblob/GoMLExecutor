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

func (d *Dictionary) getInput(n string) *shared.Field {
	if d == nil {
		return nil
	}

	input, ok := d.inputs[n]
	if !ok {
		return nil
	}

	return input
}

func (d *Dictionary) getEntry(n string) *entry {
	if d == nil {
		return nil
	}

	if len(d.registry) == 0 {
		return nil
	}

	elem, ok := d.registry[n]
	if !ok {
		return nil
	}

	if elem == nil {
		log.Printf("registry entry was nil for %v", n)
	}

	return elem
}

func (d *Dictionary) lookupString(key string, value string) (string, fieldOffset) {
	input := d.getInput(key)
	if input == nil {
		return "", unknownKeyField
	}

	ii := fieldOffset(input.Index)

	if input.Wildcard {
		return value, ii
	}

	entr := d.getEntry(key)
	if entr == nil {
		return "", unknownKeyField
	}

	if entr.hasString(value) {
		return value, ii
	}

	return oovString, ii
}

// TODO integration and boundary testing; OOV may depend on vocabulary
func (d *Dictionary) lookupInt(key string, value int) (int, fieldOffset) {
	input := d.getInput(key)
	if input 