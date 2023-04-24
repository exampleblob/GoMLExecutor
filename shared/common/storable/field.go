package storable

import (
	"github.com/viant/mly/shared/common"
	"reflect"
)

//Field represents a  default storable field descriptor
type Field struct {
	Name     string
	DataType string
	dataType reflect.Typ