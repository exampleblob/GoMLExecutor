package common

import (
	"fmt"
	"reflect"
	"strings"
)

//DataType return reflect.Type for supplied data type
func DataType(dataType string) (reflect.Type, error) {
	switch strings.ToLower(dataType) {
	case "string":
		return reflect.TypeOf(""), nil
	case "float64":
		return reflect.TypeOf(float64(0)), nil
	case "float32":
		return reflect.TypeOf(float32(0)), nil
	