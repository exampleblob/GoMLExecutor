package transfer

import (
	"fmt"
	"reflect"

	"github.com/francoispqt/gojay"
	"github.com/viant/toolbox"
)

type (
	Value interface {
		Set(values ...interface{}) error
		Key() string
		ValueAt(index int) interface{}
		UnmarshalJSONArray(dec *gojay.Decoder) error
		Len() int
		// Will make a new []interface{} that copies old value and fills the rest with the first element.
		// Panics if batchSize is less than current Values size.
		Feed(batchSize int) interface{}
	}

	Values []Value

	Strings struct {
		Name   string
		Values []string
	}
	Int32s struct {
		Name   string
		Values []int32
	}
	Int64s struct {
		Name   string
		Values []int64
	}
	Bools struct {
		Name   string
		Values []bool
	}
	Float32s struct {
		Name   string
		Values []float32
	}
	Float64s struct {
		Name   string
		Values []float64
	}
)

func (s *Strings) ValueAt(index int) interface{} {
	if index >= len(s.Values) {
		return s.Values[0]
	}
	return s.Values[index]
}

func (s *Strings) Feed(batchSize int) interface{} {
	var result = make([][]string, batchSize)
	for i, item := range s.Val