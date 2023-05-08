package common

import (
	"fmt"
	"reflect"
	"strings"
)

//DataType return reflect.Type for supplied data type
func DataType(dataType string) (reflect.Type, error) {
	s