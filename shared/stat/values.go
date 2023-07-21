package stat

import (
	"context"
	"errors"

	"github.com/viant/gmetric/stat"
)

// Values is a utility to pass multiple events to gmetric.
type Values []interface{}

func (v *Values) AppendError(err error) 