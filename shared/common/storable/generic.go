package storable

import (
	"fmt"
	"github.com/viant/mly/shared/common"
	"reflect"
)

//Generic represents generic storable
type Generic struct {
	Value interface{}
}

//Iterator returns iterator
func (s Generic) Iterator() common.Iterator {
	v := reflect.ValueOf(s.Value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var aStruct *reflectS