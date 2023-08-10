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
	for i, item := range s.Values {
		result[i] = []string{item}
	}

	for i := len(s.Values); i < batchSize; i++ {
		result[i] = []string{s.Values[0]}
		s.Values = append(s.Values, s.Values[0])
	}
	return result
}

func (s *Strings) Key() string {
	return s.Name
}

func (s *Strings) Len() int {
	return len(s.Values)
}

func (s *Strings) Set(values ...interface{}) error {
	s.Values = make([]string, len(values))
	for i, v := range values {
		s.Values[i] = toolbox.AsString(v)
	}
	return nil
}

// UnmarshalJSONArray decodes JSON array elements into slice
func (a *Strings) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var value string
	if err := dec.String(&value); err != nil {
		return err
	}
	a.Values = append(a.Values, value)
	return nil
}

// MarshalJSONArray encodes arrays into JSON
func (a Strings) MarshalJSONArray(enc *gojay.Encoder) {
	for i := 0; i < len(a.Values); i++ {
		enc.String(a.Values[i])
	}
}

// IsNil checks if array is nil
func (v Strings) IsNil() bool {
	return len(v.Values) == 0
}

func (s *Int32s) Feed(batchSize int) interface{} {
	var result = make([][]int32, batchSize)
	for i, item := range s.Values {
		result[i] = []int32{item}
	}
	for i := len(s.Values); i < batchSize; i++ {
		result[i] = []int32{s.Values[0]}
		s.Values = append(s.Values, s.Values[0])
	}
	return result
}

func (s *Int32s) ValueAt(index int) interface{} {
	if index >= len(s.Values) {
		return s.Values[0]
	}
	return s.Values[index]
}

func (s *Int32s) Key() string {
	return s.Name
}

func (s *Int32s) Len() int {
	return len(s.Values)
}

func (s *Int32s) Set(values ...interface{}) error {
	s.Values = make([]int32, len(values))
	for i, v := range values {
		val, err := toolbox.ToInt(v)
		if err != nil {
			return err
		}
		s.Values[i] = int32(val)
	}
	return nil
}

// UnmarshalJSONArray decodes JSON array elements into slice
func (a *Int32s) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var value int
	if err := dec.Int(&value); err != nil {
		return err
	}
	a.Values = append(a.Values, int32(value))
	return nil
}

// MarshalJSONArray encodes arrays into JSON
func (a Int32s) MarshalJSONArray(enc *gojay.Encoder) {
	for i := 0; i < len(a.Values); i++ {
		enc.Int(int(a.Values[i]))
	}
}

// IsNil checks if array is nil
func (v Int32s) IsNil() bool {
	return len(v.Values) == 0
}

func (s *Int64s) Feed(batchSize int) interface{} {
	var result = make([][]int64, batchSize)
	for i, ite