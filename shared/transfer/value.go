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
		// Will make a new []interface{} that copies old value and 