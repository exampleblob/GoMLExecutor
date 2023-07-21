package stat

import (
	"context"
	"errors"

	"github.com/viant/gmetric/stat"
)

// Values is a utility to pass multiple events to gmetric.
type Values []interface{}

func (v *Values) AppendError(err error) {
	if errors.Is(err, context.Canceled) {
		v.Append(Canceled)
	} else if errors.Is(err, context.DeadlineExceeded) {
		v.Append(DeadlineExceeded)
	} else {
		v.Append(err)
	}
}

func (v *Va