package storable

import (
	"fmt"
	"github.com/viant/mly/shared/common"
	"reflect"
)

//Generic represents generic storable
type Generic struct {
	Value interface{}
