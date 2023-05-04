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
	var aStruct *reflectStruct
	if v.Kind() == reflect.Struct {
		aStruct = _reflect.lookup(v.Type())
	}
	return func(pair common.Pair) error {
		switch v.Kind() {
		case reflect.Struct:
			for _, fieldType := range aStruct.fields {
				field := v.Field(fieldType.index)
				if err := pair(fieldType.name, field.Interface()); err != nil {
					return err
				}
			}
		case reflect.Map:
			switch aMap := s.Value.(type) {
			case map[string]interface{}:
				for k, v := range aMap {
					if err := pair(k, v); err != nil {
						return err
					}
				}
			case map[interface{}]interface{}:
				for k, v := range aMap {
					if err := pair(fmt.Sprintf("%s", k), v); err != nil {
						return err
					}
				}
			}
		default:
			return fmt.Errorf("unsupported generic type: %T", s.Val